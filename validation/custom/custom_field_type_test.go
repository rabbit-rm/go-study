package custom

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"

	"validation/utils"

	"github.com/go-playground/validator/v10"
)

type DbBackedUser struct {
	Name sql.NullString `validate:"required"`
	Age  sql.NullInt64  `validate:"required,gt=10"`
}

var validate *validator.Validate

func TestCustom(t *testing.T) {
	validate = validator.New()
	validate.RegisterCustomTypeFunc(validateValuer, sql.NullInt64{}, sql.NullString{})
	v := &DbBackedUser{
		Name: sql.NullString{
			String: "111",
			Valid:  true,
		},
		Age: sql.NullInt64{
			Int64: 1,
			Valid: true,
		},
	}
	utils.PrintValidatorError(validate.Struct(v))
}

func validateValuer(field reflect.Value) interface{} {
	if v, ok := field.Interface().(driver.Valuer); ok {
		value, err := v.Value()
		if err != nil {
			return nil
		}
		return value
	}
	return nil
}

type Tuple struct {
	TupleA TupleItem `validate:"required"`
	TupleB TupleItem `validate:"required,gtefield=TupleA"`
}

type TupleItem struct {
	value string
}

func (t TupleItem) Value() string {
	return t.value
}

func TestCustomStruct(t *testing.T) {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if v, ok := field.Interface().(TupleItem); ok {
			return v.value
		}
		return nil
	}, TupleItem{})
	v := Tuple{
		TupleA: TupleItem{"1"},
		TupleB: TupleItem{"2"},
	}
	utils.PrintValidatorError(validate.Struct(v))
}
