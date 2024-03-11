package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Booking struct {
	CheckIn  time.Time `form:"checkIn" binding:"required,bookableData" time_format:"2006-01-02"`
	CheckOut time.Time `form:"checkOut" binding:"required,bookableData" time_format:"2006-01-02"`
}

type Cake struct {
	Chef           Chef   `json:"chef" binding:"required"`
	CompletionTime string `json:"completionTime" binding:"required,format=2006-01-02T15:04:05"`
}

type Chef struct {
	Name         string `json:"name" binding:"required"`
	Addr         string `json:"addr" binding:"required"`
	Age          uint8  `json:"age" binding:"required"`
	WorkingYears uint8  `json:"workingYears" binding:"required,ltfield=Age"`
}

func main() {

	type S struct {
		T time.Time `json:"t"`
	}
	var s S
	err := json.Unmarshal([]byte("{\"t\":\"2024-03-04T15:04:30\"}"), &s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
	/*engine := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("format", func(fl validator.FieldLevel) bool {
			vv := fl.Field().Interface().(string)
			if _, err := time.ParseInLocation(fl.Param(), vv, time.Now().Location()); err != nil {
				return false
			}
			return true
		})
	}
	engine.GET("/cake", func(c *gin.Context) {
		var cake Cake
		if err := c.BindJSON(&cake); err != nil {
			c.JSON(400, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"status": "success",
				"cake":   cake,
			})
		}
	})
	log.Fatal(engine.Run(":80"))*/
}

func customValidator0() {
	// validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("bookableData", func(fl validator.FieldLevel) bool {
			date, ok := fl.Field().Interface().(time.Time)
			if ok {
				if time.Now().After(date) {
					return false
				}
			}
			return true
		})
	}

	e := gin.Default()
	e.GET("/validator", func(c *gin.Context) {
		var booking Booking
		if err := c.Bind(&booking); err == nil {
			c.JSON(200, gin.H{
				"status":  "success",
				"booking": booking,
			})
		} else {
			c.JSON(400, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
		}
	})
	log.Fatal(e.Run(":80"))
}
