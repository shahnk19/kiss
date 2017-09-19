package controllers

import (
	"github.com/gin-gonic/gin"
	"kiss/web/baseEnc"
	"log"
	"net/url"
)

type VM struct {
	Value  string `json:"value"`
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

func Encode() gin.HandlerFunc {
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

		encoder, err := baseEnc.Base16Encoding()
		if err != nil {
			vm.Error = err.Error()
			log.Printf("Fixme : Error encoder create. error=%v", err)
			ctx.JSON(200, vm)
			return
		}

		if encoder == nil {
			vm.Error = "Fixme! Encode is nil"
			log.Println(vm.Error)
			ctx.JSON(200, vm)
			return
		}
		vm.Value = encoder.BaseEncode(1001)
		vm.Status = true
		ctx.JSON(200, vm)
	}
}

func Decode() gin.HandlerFunc {
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

		encoder, err := baseEnc.Base16Encoding()
		if err != nil {
			vm.Error = err.Error()
			log.Printf("Fixme : Error encoder create. error=%v", err)
			ctx.JSON(200, vm)
			return
		}

		if encoder == nil {
			vm.Error = "Fixme! Encode is nil"
			log.Println(vm.Error)
			ctx.JSON(200, vm)
			return
		}
		_, err = encoder.BaseDecode("1001")
		if err != nil {
			vm.Error = err.Error()
			log.Printf("Fixme : Decoding error=%v", err)
			ctx.JSON(200, vm)
			return
		}
		vm.Status = true
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
