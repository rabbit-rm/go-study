package simple

import (
	"testing"
	"time"

	"validation/utils"

	"github.com/go-playground/validator/v10"
)

type Student struct {
	Name   string `validate:"required"`
	Email  string `validate:"email"`
	Class0 Class  `validate:"required"`
	Class1 *Class `validate:"required"`
}

type Class struct {
	No string `validate:"required"`
}
type User struct {
	FirstName      string     `validate:"required"`
	LastName       string     `validate:"required"`
	Age            uint8      `validate:"gte=0,lte=130"`
	Email          string     `validate:"required,email"`
	Gender         string     `validate:"oneof=male female prefer_not_to"`
	FavouriteColor string     `validate:"iscolor"`
	Address        []*Address `validate:"required,dive,required"`
}

type Address struct {
	Street string `validate:"required"` // 街道
	City   string `validate:"required"` // 城市
	Planet string `validate:"required"` // 行星
	Phone  string `validate:"required"` // 电话
}

var validate *validator.Validate

func TestValidator(t *testing.T) {
	validate = validator.New(validator.WithRequiredStructEnabled())
	// validateStructRequired()
	// validateStruct()
	validateVariable()
}

func validateStructRequired() {
	stu := Student{
		Name:  "rabbit",
		Email: "rabbi.rm99@gmail.com",
	}
	err := validate.Struct(stu)
	if err != nil {
		utils.PrintValidatorError(err)
	}
}

func validateStruct() {
	address := &Address{
		Street: "Rabbit Street",
		City:   "Rabbit City",
		Phone:  "none",
	}
	user := &User{
		FirstName:      "RM",
		LastName:       "Rabbit",
		Age:            34,
		Email:          "rabbit.rm99@gmail.com",
		Gender:         "male",
		FavouriteColor: "#000-",
		Address:        []*Address{address},
	}
	err := validate.Struct(user)
	if err != nil {
		utils.PrintValidatorError(err)
	}
}

func validateVariable() {
	// 自定义验证程序
	_ = validate.RegisterValidation("gtnow", func(fl validator.FieldLevel) bool {
		if v, ok := fl.Field().Interface().(time.Time); ok {
			return time.Now().After(v)
		}
		return false
	})
	email := "rabbit.rm99@gmail.com"
	wrongEmail := "12345"
	t := time.Now().Add(time.Duration(-1) * time.Second)
	if err := validate.Var(email, "required,email"); err != nil {
		utils.PrintValidatorError(err)
	}
	if err := validate.Var(wrongEmail, "required,email"); err != nil {
		utils.PrintValidatorError(err)
	}
	if err := validate.Var(t, "required,gtnow"); err != nil {
		utils.PrintValidatorError(err)
	}
}
