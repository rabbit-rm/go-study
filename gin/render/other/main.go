package main

import (
	"encoding/xml"
	"net/http"

	"gin/render/other/testdata/proto"

	"github.com/gin-gonic/gin"
)

type User struct {
	UName   string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
	Number  int64  `json:"number,omitempty"`
}

func (u User) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "userAAA",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	elem := xml.StartElement{
		Name: xml.Name{Space: "", Local: "name"},
		Attr: []xml.Attr{},
	}
	if err := e.EncodeElement(u.UName, elem); err != nil {
		return err
	}
	elem = xml.StartElement{
		Name: xml.Name{Space: "", Local: "message"},
		Attr: []xml.Attr{},
	}
	if err := e.EncodeElement(u.Message, elem); err != nil {
		return err
	}
	elem = xml.StartElement{
		Name: xml.Name{Space: "", Local: "number"},
		Attr: []xml.Attr{},
	}
	if err := e.EncodeElement(u.Number, elem); err != nil {
		return err
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func main() {
	e := gin.Default()
	e.GET("/json", func(c *gin.Context) {
		c.JSON(200, User{Number: 3456, UName: "zs", Message: "hello,rabbit.rm!"})
	})
	e.GET("/xml", func(c *gin.Context) {
		c.XML(200, User{Number: 3456, UName: "zs", Message: "hello,rabbit.rm!"})
	})
	e.GET("/yaml", func(c *gin.Context) {
		c.YAML(200, map[string]interface{}{
			"status":  http.StatusOK,
			"message": "hello,rabbit.rm!",
		})
	})
	e.GET("/proto", func(c *gin.Context) {
		label := "test"
		c.ProtoBuf(200, &proto.Test{Label: &label, Reps: []int64{1, 2, 3, 4}})
	})

	e.Run(":80")
}
