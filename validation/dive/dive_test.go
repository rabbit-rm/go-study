package dive

import (
	"testing"

	"validation/utils"

	"github.com/go-playground/validator/v10"
)

type T struct {
	Array []string          `validate:"required,gt=0,dive,required"`
	Map   map[string]string `validate:"required,gt=0,dive,keys,keymax,endkeys,required,max=1000"`
}

var validate *validator.Validate

func TestDive(t *testing.T) {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterAlias("keymax", "max=10")
	utils.PrintValidatorError(validate.Struct(T{}))
	utils.PrintValidatorError(validate.Struct(T{Array: []string{""}, Map: map[string]string{
		"test > than 10": "",
	}}))
}
