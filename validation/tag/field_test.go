package tag

import (
	"testing"

	"validation/utils"

	"github.com/go-playground/validator/v10"
)

func TestFieldTag(t *testing.T) {
	type anotherUser struct {
		No int `validate:"required" `
	}
	type user struct {
		No      int         `validate:"required,gtcsfield=Another.No"`
		Another anotherUser `validate:"required"`
		No2     int         `validate:"required,gtcsfield=No"`
		*user
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	utils.PrintValidatorError(validate.Struct(user{
		No:      3,
		No2:     4,
		Another: anotherUser{No: 1},
		user: &user{
			No:      2,
			Another: anotherUser{No: 1},
			No2:     5,
			user:    nil,
		},
	}))
}
