package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type FormA struct {
	Foo string `json:"foo" xml:"foo" form:"foo"`
}

type FormB struct {
	Bar string `json:"bar" xml:"bar" form:"bar"`
}

func bindBodyTo(ctx *gin.Context) {
	var formA FormA
	var formB FormB
	// shouldBind 读取完成Body之后不能二次读取，二次读取将EOF
	// 想要二次读取应该使用 ShouldBindWith
	if err := ctx.ShouldBind(&formA); err == nil {
		ctx.String(200, fmt.Sprintf("bind success:%v\n", formA))
	} else {
		ctx.String(500, fmt.Sprintf("error:%+v\n", err))
	}
	if err := ctx.ShouldBind(&formB); err == nil {
		ctx.String(200, fmt.Sprintf("bind success:%v", formB))
	} else {
		ctx.String(500, fmt.Sprintf("error:%+v\n\n", err))
	}
}

func bindBodyWithTo(ctx *gin.Context) {
	var formA FormA
	var formB FormB
	// shouldBind 读取完成 Body 之后不能二次读取，二次读取将EOF
	// 想要二次读取应该使用 ShouldBindBodyWith
	// 对于 Query、Form、FormPost、FormMultipart 可以多次读取
	if err := ctx.ShouldBindBodyWith(&formA, binding.JSON); err == nil {
		ctx.String(200, fmt.Sprintf("bind success:%v\n", formA))
	} else {
		ctx.String(500, fmt.Sprintf("error:%+v\n", err))
	}
	if err := ctx.ShouldBindBodyWith(&formB, binding.JSON); err == nil {
		ctx.String(200, fmt.Sprintf("bind success:%v\n", formB))
	} else {
		ctx.String(500, fmt.Sprintf("error:%+v\n", err))
	}
}
