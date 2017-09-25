package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mediocregopher/radix.v2/pool"
	"kiss/web/baseEnc"
	"kiss/web/models"
	"log"
	"net/url"
)

const (
	ERR_NEW_ENTRY_FAIL        = "NewEntry failed, Error=%v"
	ERR_NIL_ENCODER           = "Fixme! Encoder is nil"
	ERR_INVALID_URL           = "Encode: Invalid url provided, url=%s"
	ERR_DESTINATION_NOT_FOUND = "Destination url for code(%s),id(%d) not found.Error=%v"
)

type Ctrl struct {
	model   models.IModel
	encoder *baseEnc.Encoding
	cache   *pool.Pool
	Name    string
}

func New(conn string) *Ctrl {
	p, err := pool.New("tcp", "localhost:6379", 500)
	if err != nil {
		log.Panic(err)
	}
	return &Ctrl{
		Name:    "hotpie",
		encoder: getBaseEncoder(),
		model:   models.New(conn),
		cache:   p,
	}
}

type VM struct {
	Value  string `json:"value"`
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

func Encode(c *Ctrl) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		vm := &VM{
			Value:  "",
			Status: false,
			Error:  "",
		}

		url := ctx.Request.FormValue("url")
		if !isValidUrl(url) {
			vm.Error = fmt.Sprintf(ERR_INVALID_URL, url)
			log.Printf(vm.Error)
			ctx.JSON(200, vm)
			return
		}

		if c.encoder == nil {
			vm.Error = ERR_NIL_ENCODER
			log.Println(vm.Error)
			ctx.JSON(200, vm)
			return
		}
		if id, err := c.model.GetIdByUrl(url); err == nil && id > 0 {
			vm.Value = c.encoder.BaseEncode(id)
			vm.Status = true
			ctx.JSON(200, vm)
			return
		}
		tinyUrl, err := c.model.SaveTiny(url, c.encoder)
		if err != nil {
			vm.Error = fmt.Sprintf(ERR_NEW_ENTRY_FAIL, err)
			log.Println(vm.Error)
			ctx.JSON(200, vm)
			return
		}

		vm.Value = tinyUrl
		vm.Status = true
		ctx.JSON(200, vm)
	}
}

func Decode(c *Ctrl) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		vm := &VM{
			Value:  "",
			Status: false,
			Error:  "",
		}

		code := ctx.Request.FormValue("code")

		if c.encoder == nil {
			vm.Error = ERR_NIL_ENCODER
			log.Println(vm.Error)
			ctx.JSON(200, vm)
			return
		}
		id, err := c.encoder.BaseDecode(code)
		if err != nil {
			vm.Error = err.Error()
			log.Printf("Fixme : Decoding error=%v", err)
			ctx.JSON(200, vm)
			return
		}
		//load from db rows with id
		url, err := c.model.GetUrlById(id)
		if err != nil {
			vm.Error = fmt.Sprintf(ERR_DESTINATION_NOT_FOUND, code, id, err)
			log.Printf(vm.Error)
			ctx.JSON(200, vm)
			return
		}

		vm.Value = url
		ctx.JSON(200, vm)
	}
}

func Parser(c *Ctrl) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		cacheConn, err := c.cache.Get()
		if err != nil {
			log.Panicln(err)
		}

		code := ctx.Request.URL.Path
		code = code[1:] //remove the first chat '/'
		url := ""
		//try load from cache
		if cacheConn != nil {
			url, err = cacheConn.Cmd("HGET", code, "url").Str()
			if err != nil {
				log.Printf("Error getting from cache, error = %v", err)
			} else {
				c.cache.Put(cacheConn)
				ctx.Redirect(302, url)
				return
			}
		}

		if c.encoder == nil {
			log.Println(ERR_NIL_ENCODER)
			c.cache.Put(cacheConn)
			ctx.String(404, "404 Not Found")
			return
		}
		id, err := c.encoder.BaseDecode(code)
		if err != nil {
			log.Printf("Fixme : Decoding error=%v", err)
			c.cache.Put(cacheConn)
			ctx.String(404, "404 Not Found")
			return
		}
		//load from db rows with id
		url, err = c.model.GetUrlById(id)
		if err != nil {
			log.Printf(fmt.Sprintf(ERR_DESTINATION_NOT_FOUND, code, id, err))
			c.cache.Put(cacheConn)
			ctx.String(404, "404 Not Found")
			return
		}
		if url == "" {
			c.cache.Put(cacheConn)
			ctx.String(404, "404 Not Found")
			return
		}

		//save to redis so next time we don't come here again
		if cacheConn != nil {
			err = cacheConn.Cmd("HMSET", code, "url", url).Err
			if err != nil {
				log.Printf("Error setting to cache, error = %v", err)
			}
		}
		c.cache.Put(cacheConn)
		ctx.Redirect(302, url)
	}
}

func isValidUrl(sUrl string) bool {
	if sUrl == "" {
		return false
	}
	_, err := url.ParseRequestURI(sUrl)
	if err != nil {
		return false
	}
	return true
}

func getBaseEncoder() *baseEnc.Encoding {
	encoder, err := baseEnc.Base16Encoding()
	if err != nil {
		log.Panic("Fixme : Error encoder create. error=%v", err)
	}
	return encoder
}
