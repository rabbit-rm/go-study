package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func PrintValidatorError(err error) {
	if err == nil {
		return
	}
	var invalidValidationError *validator.InvalidValidationError
	if errors.As(err, &invalidValidationError) {
		fmt.Println(err)
		return
	}
	var v validator.ValidationErrors
	if errors.As(err, &v) {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Error())
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		return
	}
	fmt.Println(err)
}
