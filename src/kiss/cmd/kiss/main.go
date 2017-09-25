package main

import (
	"flag"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"html/template"
	"kiss/web/controllers"
	"log"
	"os"
	"strings"
	"syscall"
)

var (
	port int
)

func init() {
	flag.IntVar(&port, "port", 5000, "web server port number")
}

func parseTemplates() multitemplate.Render {
	r := multitemplate.New()

	withHeaderAndFooter := []string{
		"index.html",
	}

	tdir := "templ/"
	for _, file := range withHeaderAndFooter {
		t := template.New(file)
		t.Funcs(template.FuncMap{"appName": func() string {
			return strings.TrimSuffix(t.Name(), ".html") + "App"
		}})
		r.Add(file, template.Must(t.ParseFiles(tdir+file)))
	}
	return r
}

func main() {
	flag.Parse()

	route := gin.Default()
	route.Use(static.Serve("/w/css", static.LocalFile("css", true)))

	route.HTMLRender = parseTemplates()
	route.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	route.GET("/w/index", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{
			"version": 1,
		})
	})
	if ctrl := controllers.New("postgres://postgres:root@localhost/kiss?sslmode=disable"); ctrl != nil {
		route.NoRoute(controllers.Parser(ctrl))
		apiRouteGroup := route.Group("/api")
		{
			apiRouteGroup.GET("/encode", controllers.Encode(ctrl))
			apiRouteGroup.GET("/decode", controllers.Decode(ctrl))
		}
	}

	listenAddr := fmt.Sprintf("localhost:%d", port)
	graceful := endless.NewServer(listenAddr, route)
	err := graceful.RegisterSignalHook(endless.PRE_SIGNAL, syscall.SIGHUP, func() {
		log.Println("SIGHUP - Shutting down webserver")
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := graceful.ListenAndServe(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("shutting down")
	}
	os.Exit(0)
}
