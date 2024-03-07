package main

import (
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
)

var tmpl = template.Must(template.New("https").Parse(`
<html>
<head>
  <title>Https Test</title>
  <script src="/assets/app.js"></script>
</head>
<body>
  <h1 style="color:red;">Welcome, Rabbit.RM!</h1>
</body>
</html>
`))

func main() {
	// 服务器推送静态资源
	e := gin.Default()
	e.Static("/assets", "W:\\GoProject\\private\\gin\\serverPush\\assets")
	e.SetHTMLTemplate(tmpl)
	e.GET("/index", func(c *gin.Context) {
		if pusher := c.Writer.Pusher(); pusher != nil {
			if err := pusher.Push("/assets/app.js", nil); err != nil {
				log.Printf("ERROR:%+v", err)
			}
		}
		c.HTML(200, "https", gin.H{
			"status": "success",
		})
	})
	certFile := "W:\\GoProject\\private\\gin\\serverPush\\certificate\\cert.pem"
	keyFile := "W:\\GoProject\\private\\gin\\serverPush\\certificate\\key.pem"
	_ = e.RunTLS(":80", certFile, keyFile)
}
