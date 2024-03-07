package main

import (
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	renderHtml3()
}

func renderHtml() {
	e := gin.Default()
	e.LoadHTMLGlob("W:\\GoProject\\private\\gin\\html\\templates\\*")
	e.GET("/index", func(context *gin.Context) {
		context.HTML(200, "index.tmpl", gin.H{
			"title": "Rabbit.rm",
		})
	})
	_ = e.Run(":80")
}

// 响应不同目录下的模板
func renderHtml2() {
	e := gin.Default()
	e.LoadHTMLGlob("W:\\GoProject\\private\\gin\\html\\templates\\**\\*")
	e.GET("/user/index", func(c *gin.Context) {
		c.HTML(200, "user/index.tmpl", gin.H{
			"title": "Rabbit.rm(user/index)",
		})
	})
	e.GET("/index/index", func(c *gin.Context) {
		c.HTML(200, "index/index.tmpl", gin.H{
			"title": "Rabbit.rm(index/index)",
		})
	})
	e.Run(":80")
}

// 自定义模板渲染器
func renderHtml3() {
	e := gin.Default()
	// 更改模板分隔符
	e.Delims("{{{", "}}}")
	// 2种方式都行,内部调用都一样
	/*e.SetFuncMap(map[string]any{
		"fmtBirthday": func(date time.Time) string {
			return date.Format("2006-01-02T15:04:05.000")
		},
	})
	e.LoadHTMLFiles("W:\\\\GoProject\\\\private\\\\gin\\\\html\\\\templates\\\\user\\\\user.template")*/
	templ := template.Must(template.New("").Delims("{{{", "}}}").Funcs(map[string]any{
		"fmtBirthday": func(date time.Time) string {
			return date.Format("2006-01-02T15:04:05.000")
		},
	}).ParseFiles("W:\\GoProject\\private\\gin\\html\\templates\\user\\user.template"))
	e.SetHTMLTemplate(templ)
	e.GET("/users", func(c *gin.Context) {
		c.HTML(200, "user/user.template", gin.H{
			"users": []struct {
				Name     string
				Age      uint8
				Birthday time.Time
			}{
				{Name: "zs", Age: 34, Birthday: time.Date(1999, 3, 5, 20, 40, 50, 0, time.Now().Location())},
				{Name: "ls", Age: 24, Birthday: time.Date(1998, 3, 5, 20, 40, 50, 0, time.Now().Location())},
				{Name: "ww", Age: 36, Birthday: time.Date(2001, 3, 5, 20, 40, 50, 0, time.Now().Location())},
				{Name: "rabbit", Age: 14, Birthday: time.Date(2003, 3, 5, 20, 40, 50, 0, time.Now().Location())},
			},
		})

	})
	e.Run(":80")
}
