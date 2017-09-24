package main

import (
	"flag"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"kiss/web/controllers"
	"log"
	"os"
	"syscall"
)

var (
	port int
)

func init() {
	flag.IntVar(&port, "port", 5000, "web server port number")
}

func main() {
	flag.Parse()
	route := gin.Default()
	route.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	if ctrl := controllers.New("postgres://postgres:root@localhost/kiss"); ctrl != nil {
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
