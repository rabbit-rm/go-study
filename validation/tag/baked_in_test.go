package tag

import (
	"testing"

	"validation/utils"

	"github.com/go-playground/validator/v10"
)

func TestStructOnly(t *testing.T) {
	type Address struct {
		Street string `validate:"required"`
		Phone  string `validate:"required"`
	}
	type User struct {
		Name    string   `validate:"required"`
		Address *Address `validate:"required,nostructlevel"`
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	utils.PrintValidatorErrorSimplify(validate.Struct(&User{
		Name:    "rabbit",
		Address: &Address{},
	}))
}

func TestOmitempty(t *testing.T) {
	type Number struct {
		Num1 int `validate:"required,omitempty"`
		Num2 int `validate:"omitempty,max=10,min=1"`
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	utils.PrintValidatorErrorSimplify(validate.Struct(&Number{
		Num1: 20,
		Num2: 11,
	}))
}

func TestDive(t *testing.T) {
	type T struct {
		Slice []int64  `validate:"required,dive,max=10"`
		Array [5]int64 `validate:"required,dive,min=20"`
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Slice: []int64{1, 2, 3, 14, 5},
		Array: [5]int64{20, 21, 2, 23, 34},
	}))

	type TT struct {
		// 跳过第二层
		Array2 [][]int64 `validate:"required,gt=3,dive,dive,required,max=5"`
	}
	utils.PrintValidatorErrorSimplify(validate.Struct(&TT{
		Array2: [][]int64{
			{1},
			{1, 2},
			{1, 2, 3},
			{1, 2, 3, 4},
			{1, 2, 3, 4, 5},
		},
	}))

}

func TestDiveMap(t *testing.T) {

	type TTM struct {
		// required,gt=3：作用于 map[int64]string，验证 map 长度必须大于3
		// dive,keys,max=10,min=5,endkeys,lt=10：作用于 map 子项，表示验证 key在5-10之间，value 长度小于10
		Map map[int64]string `validate:"required,gt=3,dive,keys,max=10,min=5,endkeys,lt=10"`
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	utils.PrintValidatorErrorSimplify(validate.Struct(&TTM{
		Map: map[int64]string{
			6: "rabbit.rm",
			7: "rabbit.rm",
			8: "rabbit.rm,rabbit.rm,rabbit.rm",
			9: "rabbit.rm",
		},
	}))

	type TTM2 struct {
		// required,gt=3：作用于 map[int64]map[int64]string，验证 map 长度必须大于3
		// dive,keys,max=10,min=5,endkeys,len=5：作用于 map 子项，表示验证 key 在5-10之间，value（map[int64]string） 长度等于5
		// dive,keys,max=5,min=0,endkeys,lt=10：作用于 map[int64]string，表示验证 表示验证 key 在0-5之间，value（string） 长度小于10
		Map map[int64]map[int64]string `validate:"required,gt=3,dive,keys,max=10,min=5,endkeys,len=5,dive,keys,max=5,min=0,endkeys,lt=10"`
	}
	utils.PrintValidatorErrorSimplify(validate.Struct(&TTM2{map[int64]map[int64]string{
		6: {0: "rabbit.rm,rabbit.rm", 1: "", 2: "", 3: "", 4: ""},
		7: {0: "", 1: "", 2: "", 3: "", 4: ""},
		8: {0: "", 1: "", 2: "", 3: "", 4: ""},
		9: {0: "", 1: "", 2: "", 3: "", 4: ""},
	}}))
}

func TestRequiredIf(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	type T struct {
		Enable1 bool
		Enable2 bool
		// Enable1=true
		Port1 uint64 `validate:"required_if=Enable1 true"`
		// Enable1=true && Enable2=true
		Port2 uint64 `validate:"required_if=Enable1 true Enable2 true"`
	}
	// PASS
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: false,
		Port1:   0,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: false,
		Port1:   1,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: true,
		Port1:   0,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: true,
		Port1:   1,
		Port2:   1,
	}))
}

func TestRequiredUnless(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	type T struct {
		Enable1 bool
		Enable2 bool
		// Enable1 == false
		Port1 uint64 `validate:"required_unless=Enable1 true"`
		// Enable1 == false && enable2 == false
		Port2 uint64 `validate:"required_unless=Enable1 true Enable2 true"`
	}
	// PASS
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: false,
		Port1:   1,
		Port2:   1,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: false,
		Port1:   0,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: true,
		Port1:   1,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: true,
		Port1:   0,
		Port2:   0,
	}))
}

func TestRequiredWith(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	type T struct {
		Enable1 bool
		Enable2 bool
		// Enable1非零值
		Port1 uint64 `validate:"required_with=Enable1"`
		// Enable1非零值 || enable2非零值
		Port2 uint64 `validate:"required_with=Enable1 Enable2"`
	}
	// PASS
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: false,
		Port1:   0,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: false,
		Port1:   1,
		Port2:   1,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: true,
		Port1:   0,
		Port2:   1,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: true,
		Port1:   1,
		Port2:   1,
	}))
}

func TestRequiredWithAll(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	type T struct {
		Enable1 bool
		Enable2 bool
		// Enable1非零值
		Port1 uint64 `validate:"required_with_all=Enable1"`
		// Enable1非零值 && enable2非零值
		Port2 uint64 `validate:"required_with_all=Enable1 Enable2"`
	}
	// PASS
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: false,
		Port1:   0,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: false,
		Port1:   1,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: true,
		Port1:   0,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: true,
		Port1:   1,
		Port2:   1,
	}))
}

func TestRequiredWithout(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	type T struct {
		Enable1 bool
		Enable2 bool
		// Enable1非零验证该字段必须有值，其他忽略
		Port1 uint64 /*`validate:"required_without=Enable1"`*/
		// Enable1非零||Enable2非零，验证该字段必须有值，其他忽略
		Port2 uint64 `validate:"required_without=Enable1 Enable2"`
	}
	// PASS
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: false,
		Port1:   1,
		Port2:   1,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: false,
		Port1:   0,
		Port2:   1,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: true,
		Port1:   1,
		Port2:   1,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: true,
		Port1:   0,
		Port2:   1,
	}))
}

func TestRequiredWithoutAll(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	type T struct {
		Enable1 bool
		Enable2 bool
		// Enable1是零值
		Port1 uint64 `validate:"required_without_all=Enable1"`
		// Enable1是零值 && enable2是零值
		Port2 uint64 `validate:"required_without_all=Enable1 Enable2"`
	}
	// PASS
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: false,
		Port1:   1,
		Port2:   1,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: false,
		Port1:   0,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: true,
		Port1:   1,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: true,
		Port1:   0,
		Port2:   0,
	}))
}

func TestExcludedIf(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	type T struct {
		Enable1 bool
		Enable2 bool
		// Enable1 == true 验证是否为零值，其他情况忽略
		Port1 uint64 `validate:"excluded_if=Enable1 true"`
		// Enable1 == true && Enable2 == true 验证是否为零值，其他情况忽略
		Port2 uint64 `validate:"excluded_if=Enable1 true Enable2 true"`
	}
	// PASS
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: false,
		Port1:   1,
		Port2:   1,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: false,
		Port1:   0,
		Port2:   1,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: true,
		Port1:   1,
		Port2:   1,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: true,
		Port1:   0,
		Port2:   0,
	}))
}

func TestExcludedUnless(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	type T struct {
		Enable1 bool
		Enable2 bool
		// Enable1 != true 验证是否为零值，其他情况忽略
		Port1 uint64 `validate:"excluded_unless=Enable1 true"`
		// Enable1 != true || Enable2 != true 验证是否为零值，其他情况忽略
		Port2 uint64 `validate:"excluded_unless=Enable1 true Enable2 true"`
	}
	// PASS
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: false,
		Port1:   0,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: false,
		Port1:   1,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: false,
		Enable2: true,
		Port1:   0,
		Port2:   0,
	}))
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{
		Enable1: true,
		Enable2: true,
		Port1:   1,
		Port2:   1,
	}))
}

func TestIsDefault(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	type T struct {
		Port uint64 `validate:"isdefault=1"`
	}
	// PASS
	utils.PrintValidatorErrorSimplify(validate.Struct(&T{Port: 1}))

}
