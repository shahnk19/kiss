package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kiss/web/baseEnc"
	"kiss/web/models"
	"log"
	"net/url"
)

const (
	ERR_NEW_ENTRY_FAIL        = "NewEntry failed"
	ERR_NIL_ENCODER           = "Fixme! Encoder is nil"
	ERR_INVALID_URL           = "Encode: Invalid url provided, url=%s"
	ERR_DESTINATION_NOT_FOUND = "Destination url for code(%s),id(%d) not found.Error=%v"
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

		retryCounter := 5
		lastId := c.encoder.BaseEncode(c.model.GetLastId())
		var err error
		var tinyUrl string
		//insert into db, and get the new row id
		for retryCounter >= 0 {
			tinyUrl, err = c.model.SaveTiny(url, lastId)
			if err == nil && tinyUrl != "" {
				break
			}
			retryCounter--
		}

		if retryCounter <= 0 { //can't get anything unique
			vm.Error = ERR_NEW_ENTRY_FAIL
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
		url, err := c.model.GetDestination(id)
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
