package mapValidation

import (
	"fmt"
	"strings"
	"testing"

	"validation/utils"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func TestMapValidate(t *testing.T) {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validateMap()

}

func validateMap() {
	user := map[string]interface{}{
		"name":  "1rabbit.RM",
		"email": "rabbit.rm99@gmail.com",
	}
	rules := map[string]interface{}{
		"name":  "required,isPrefix=rabbit",
		"email": "required,email",
	}
	_ = validate.RegisterValidation("isPrefix", func(fl validator.FieldLevel) bool {
		if v, ok := fl.Field().Interface().(string); ok {
			param := fl.Param()
			return strings.HasPrefix(v, param)
		}
		return false
	})
	m := validate.ValidateMap(user, rules)
	for k, v := range m {
		if vv, ok := v.(error); ok {
			fmt.Println(k)
			utils.PrintValidatorError(vv)
		}
	}
}
