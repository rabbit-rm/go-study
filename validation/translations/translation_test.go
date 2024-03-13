package translations

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func TestTranslation(t *testing.T) {
	zhT := zh.New()
	uni = ut.New(zhT, zhT)

	trans, found := uni.GetTranslator("zh")
	if !found {
		log.Fatal("not found en translation")
	}
	validate = validator.New(validator.WithRequiredStructEnabled())
	_ = zhtranslations.RegisterDefaultTranslations(validate, trans)
	translateAll(trans)
	translateIndividual(trans)
	translateOverride(trans)
}

func translateAll(trans ut.Translator) {

	type User struct {
		Username string `validate:"required"`
		Tagline  string `validate:"required,lt=10"`
		Tagline2 string `validate:"required,gt=1"`
	}

	user := User{
		Username: "rabbit",
		Tagline:  "14",
		Tagline2: "2",
	}

	err := validate.Struct(user)
	if err != nil {
		// 翻译所有
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			fmt.Println(errs.Translate(trans))
		}
	}
}

// 翻译单个
func translateIndividual(trans ut.Translator) {
	type User struct {
		Username string `validate:"required"`
		Tagline  string `validate:"required,lt=10"`
		Tagline2 string `validate:"required,gt=1"`
	}

	user := User{
		Username: "rabbit",
		Tagline:  "14",
		Tagline2: "11",
	}

	err := validate.Struct(user)
	if err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			for _, err := range errs {
				fmt.Println(err.Translate(trans))
			}
		}
	}
}

func translateOverride(trans ut.Translator) {
	_ = validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0}必填", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	type User struct {
		Username string `validate:"required"`
	}

	var user User

	err := validate.Struct(user)
	if err != nil {

		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			for _, e := range errs {
				fmt.Println(e.Translate(trans))
			}
		}
	}

}
