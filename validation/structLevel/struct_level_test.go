package structLevel

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"validation/utils"

	"github.com/go-playground/validator/v10"
)

type validationError struct {
	Namespace       string `json:"namespace"`
	Field           string `json:"field"`
	StructNamespace string `json:"structNamespace"`
	StructField     string `json:"structField"`
	Tag             string `json:"tag"`
	ActualTag       string `json:"actualTag"`
	Kind            string `json:"kind"`
	Type            string `json:"type"`
	Value           string `json:"value"`
	Param           string `json:"param"`
	Message         string `json:"message"`
}

type gender uint8

const (
	Male gender = iota + 1
	Female
	Intersex
)

func (g gender) String() string {
	var terms = [3]string{"Male", "Female", "Intersex"}
	if g < Male || g > Intersex {
		return "UNKNOWN"
	}
	return terms[g]
}

// user information
type User struct {
	FirstName      string     `json:"fname"`
	LastName       string     `json:"lname"`
	Age            uint8      `validate:"gte=0,lte=130"`
	Email          string     `json:"e-mail" validate:"required,email"`
	FavouriteColor string     `validate:"hexcolor|rgb|rgba"`
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
	Gender         gender     `json:"gender" validate:"required,gender_custom_validation"`
}

// 地址包含用户地址信息
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

func TestStructLevel(t *testing.T) {
	var validate *validator.Validate
	validate = validator.New()
	// 注册函数以从json标签获取标签名称
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := field.Tag.Get("json")
		name := strings.SplitN(tag, ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// 为 User 注册验证
	// 注意: 只需要为 User 注册一个非指针类型验证器
	// 在它的类型检查过程中内部解引用
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		user := sl.Current().Interface().(User)
		if len(user.FirstName) == 0 && len(user.LastName) == 0 {
			sl.ReportError(user.FirstName, "fname", "FirstName", "fnameorlname", "")
			sl.ReportError(user.LastName, "lname", "LastName", "fnameorlname", "")
		}

	}, User{})

	// 在一行上注册用户类型的自定义验证函数
	// 验证枚举是否在间隔内
	if err := validate.RegisterValidation("gender_custom_validation", func(fl validator.FieldLevel) bool {
		if v, ok := fl.Field().Interface().(gender); ok {
			return v.String() != "UNKNOWN"
		}
		return false
	}); err != nil {
		utils.PrintValidatorError(err)
	}

	address := &Address{
		Street: "Rabbit Street",
		Planet: "Rabbit Planet",
		Phone:  "none",
		City:   "Unknown",
	}

	user := &User{
		FirstName:      "",
		LastName:       "",
		Age:            45,
		Email:          "rabbit.rm99gmail.com",
		FavouriteColor: "#000",
		Addresses:      []*Address{address},
		Gender:         Male,
	}

	if err := validate.Struct(user); err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			utils.PrintValidatorError(err)
			return
		}
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, err := range err.(validator.ValidationErrors) {
				e := validationError{
					Namespace:       err.Namespace(),
					Field:           err.Field(),
					StructNamespace: err.StructNamespace(),
					StructField:     err.StructField(),
					Tag:             err.Tag(),
					ActualTag:       err.ActualTag(),
					Kind:            fmt.Sprintf("%v", err.Kind()),
					Type:            fmt.Sprintf("%v", err.Type()),
					Value:           fmt.Sprintf("%v", err.Value()),
					Param:           err.Param(),
					Message:         err.Error(),
				}
				indent, err := json.MarshalIndent(e, "", "  ")
				if err != nil {
					utils.PrintValidatorError(err)
					return
				}
				fmt.Println(string(indent))
			}
		}
	}
}
