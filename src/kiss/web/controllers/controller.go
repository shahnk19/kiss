package controllers

import (
	"github.com/gin-gonic/gin"
	"kiss/web/baseEnc"
	"kiss/web/models"
	"log"
	"net/url"
)

type Ctrl struct {
	model   models.IModel
	encoder *baseEnc.Encoding
	Name    string
}

func New(conn string) *Ctrl {
	return &Ctrl{
		Name:    "hotpie",
		encoder: getBaseEncoder(),
		model:   models.New(conn),
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
			vm.Error = "Invalid url"
			log.Printf("Encode: Invalid url provided, url=%s", url)
			ctx.JSON(200, vm)
			return
		}

		if c.encoder == nil {
			vm.Error = "Fixme! Encoder is nil"
			log.Println(vm.Error)
			ctx.JSON(200, vm)
			return
		}

		//insert into db, and get the new row id
		newId, err := c.model.NewEntry(url)
		if err != nil {
			vm.Error = "NewEntry failed"
			log.Println(vm.Error + err.Error())
			ctx.JSON(200, vm)
		}

		vm.Value = c.encoder.BaseEncode(newId)
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
			vm.Error = "Fixme! Encoder is nil"
			log.Println(vm.Error)
			ctx.JSON(200, vm)
			return
		}
		_, err := c.encoder.BaseDecode(code)
		if err != nil {
			vm.Error = err.Error()
			log.Printf("Fixme : Decoding error=%v", err)
			ctx.JSON(200, vm)
			return
		}
		vm.Status = true

		//load from db rows with id = vid

		vm.Value = ""
		ctx.JSON(200, vm)
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
