package structLevel

import (
	"fmt"
	"testing"

	"validation/utils"

	"github.com/go-playground/validator/v10"
)

type HighSchoolStudent struct {
	Name  string `validate:"required"`
	Age   uint8  `validate:"required,minAge=15"`
	Clazz Class  `validate:"required"`
}

type Class struct {
	No   string `validate:"required,checkClassNo"`
	Name string `validate:"required"`
}

func TestFieldLevel(t *testing.T) {
	var a = uint8(1)
	var age = &a
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())
	_ = validate.RegisterValidation("minAge", func(fl validator.FieldLevel) bool {
		printFieldLevel(fl)
		return true
	})
	_ = validate.RegisterValidation("checkClassNo", func(fl validator.FieldLevel) bool {
		printFieldLevel(fl)
		return true
	})
	utils.PrintValidatorError(validate.Struct(&HighSchoolStudent{
		Name: "rabbit",
		Age:  10,
		Clazz: Class{
			No:   "No:123",
			Name: "ROOM-2",
		},
	}))
	utils.PrintValidatorError(validate.Var(age, "minAge"))
}

func printFieldLevel(fl validator.FieldLevel) {
	fmt.Printf("Top:%v\n", fl.Top())
	fmt.Printf("Parent:%v\n", fl.Parent())
	fmt.Printf("Field:%v\n", fl.Field())
	fmt.Printf("FieldName:%v\n", fl.FieldName())
	fmt.Printf("StructFieldName:%v\n", fl.StructFieldName())
	fmt.Printf("Param:%v\n", fl.Param())
	fmt.Printf("GetTag:%v\n", fl.GetTag())
	value, kind, nullable := fl.ExtractType(fl.Field())
	fmt.Printf("ExtractType:value[%v],kind:[%v],nullable[%v]\n", value, kind, nullable)
	// fmt.Printf("GetStructFieldOK2:%v", fl.GetStructFieldOK2())
	// fmt.Printf("GetStructFieldOKAdvanced2:%v", fl.GetStructFieldOKAdvanced2())
	fmt.Println("______________________________________________")
}
