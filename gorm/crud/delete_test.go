package crud

import (
	"fmt"
	"testing"

	"gorm.io/gorm"
)

func TestDeleteMethod(t *testing.T) {
	// 指定主键删除
	printResult(db.Model(&User{}).Delete(&User{}, "id = ?", 15))
	// 批量删除
	printResult(db.Model(&User{}).Where("name LIKE ?", "user-_").Delete(&User{}))
	// 不带 Where 子句删除,ErrMissingWhereClause
	printResult(db.Model(&User{}).Delete(&User{}))

	printResult(db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&User{}))
}

func TestUnscopedDelete(t *testing.T) {
	var users []User
	printResult(db.Model(&User{}).Where("deleted_at IS NOT NULL").Unscoped().Find(&users))
	fmt.Printf("users:%v\n", users)
}

func TestSoftDeleteUnixTimestamp(t *testing.T) {
	/*for i := 0; i < 10; i++ {
		db.Model(&softDeleteUnixTimestamp{}).Create(&softDeleteUnixTimestamp{
			Name: fmt.Sprintf("name_%d", i),
		})
	}*/
	db.Model(&softDeleteUnixTimestamp{}).Delete(nil, []int{2, 4, 6, 8, 10})
	var dest []softDeleteUnixTimestamp
	db.Model(&softDeleteUnixTimestamp{}).Unscoped().Find(&dest, []int{2, 4, 6, 8, 10})
	fmt.Printf("dest:%v\n", dest)
}

func TestSoftDeleteFlag(t *testing.T) {
	for i := 0; i < 10; i++ {
		db.Model(&softDeleteFlag{}).Create(&softDeleteFlag{
			Name: fmt.Sprintf("name_%d", i),
		})
	}
}

func TestSoftDeleteMixed(t *testing.T) {
	for i := 0; i < 10; i++ {
		db.Model(&softDeleteMixed{}).Create(&softDeleteMixed{
			Name: fmt.Sprintf("name_%d", i),
		})
	}
}
