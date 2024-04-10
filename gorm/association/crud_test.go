package association

import (
	"fmt"
	"testing"

	"gorm.io/gorm/clause"
)

func TestFindUser(t *testing.T) {
	var user User
	printResult(db.Model(&User{}).Where("id = ?", 1).Find(&user))
	fmt.Printf("user:%#v\n", user)
	// preload company
	printResult(db.Model(&User{}).Preload("Company").Where("id = ?", 1).Find(&user))
	fmt.Printf("user:%#v\n", user.Company.Name)

	// preload identity card
	printResult(db.Model(&User{}).Preload("IdentityCard").Where("id = ?", 1).Find(&user))
	fmt.Printf("user.IdentityCard.No:%#v,user.IdentityCard.Address:%#v\n", user.IdentityCard.No, user.IdentityCard.Address)

	// preload identity card
	printResult(db.Model(&User{}).Preload("Courses").Where("id = ?", 1).Find(&user))
	fmt.Printf("courses:%#v\n", user.Courses)

	// 关联查询使用 association
	// association course
	var courses []Course
	err := db.Model(&user).Association("Courses").Find(&courses)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
	}
}

func TestAssociationSelectFind(t *testing.T) {
	var company Company
	var user User
	tx := db.Model(&User{}).Preload("Company").Preload("IdentityCard").Find(&user, "id = ?", 19)
	printResult(tx)
	fmt.Printf("user:%#v\n", &user)
	err := db.Model(&user).Select("name").Association("Company").Find(&company)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
	}
	fmt.Printf("company:%#v\n", &company)
}

func TestAssociationDelete(t *testing.T) {
	{
		// 仅删除用户，不删除关联关系
		var user User
		tx := db.Model(&User{}).Find(&user, "id = ?", 17)
		printResult(tx)
		fmt.Printf("user:%#v\n", &user)
		tx = db.Delete(&user)
		printResult(tx)
	}
	{
		// 删除用户，并删除 IdentityCard 关联
		// delete user IdentityCard
		var user User
		tx := db.Model(&User{}).Preload("Company").Preload("IdentityCard").Find(&user, "id = ?", 18)
		printResult(tx)
		fmt.Printf("user:%#v\n", &user)
		tx = db.Select("IdentityCard").Delete(&user)
		printResult(tx)
	}
	{
		// 删除用户，删除所有关联关系，软删除 IdentityCard 移除 user-course 关系表数据，删除用户
		var user User
		tx := db.Model(&User{}).Preload("IdentityCard").Preload("Company").Find(&user, "id = ?", 15)
		printResult(tx)
		fmt.Printf("user:%#v\n", &user)
		tx = db.Select(clause.Associations).Delete(&user)
		printResult(tx)
	}
}

func TestAssociationFind(t *testing.T) {
	// 查询关联
	var user User
	db.Model(&User{}).Find(&user, "id = ?", 20)
	// 简单 关联查询
	err := db.Model(&user).Association("Company").Find(&user.Company)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
	}
	fmt.Printf("user:%#v\n", &user.Company)
	// 条件关联查询
	err = db.Model(&user).Where("name = ?", "Golang").Association("Courses").Find(&user.Courses)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
	}
	for _, course := range user.Courses {
		fmt.Printf("user:%#v\n", &course)
	}
}

func TestAssociationAppendOen2One(t *testing.T) {
	{
		var user User
		db.Model(&User{}).Preload("Company").Find(&user, "id = ?", 21)
		fmt.Printf("user:%#v\n", &user)
		var companys []Company
		db.Model(&Company{}).Find(&companys)
		fmt.Printf("company:%#v\n", &companys)
		// 一对一关系，直接替换原来的关联
		err := db.Model(&user).Association("Company").Append(&companys[0])
		if err != nil {
			fmt.Printf("err:%#v\n", err)
		}
	}

	{
		var user User
		db.Model(&User{}).Preload("Courses").Find(&user, "id = ?", 22)
		fmt.Printf("user:%#v\n", &user)
		var courses []Course
		db.Model(&Course{}).Find(&courses)
		for _, course := range courses {
			fmt.Printf("course:%#v\n", &course)
		}
		// 多对多关系，追加关联，会出现重复数据
		err := db.Model(&user).Association("Courses").Append(&courses[0], &courses[1])
		if err != nil {
			fmt.Printf("err:%#v\n", err)
		}
		fmt.Printf("user:%#v\n", &user)
	}

}

func TestAssociationReplace(t *testing.T) {
	{
		var user User
		db.Model(&User{}).Preload("Company").Find(&user, "id = ?", 23)
		fmt.Printf("user:%#v\n", &user)
		var company Company
		db.Model(&Company{}).Find(&company, "id = ?", 4)
		err := db.Model(&user).Association("Company").Replace(&company)
		if err != nil {
			fmt.Printf("err:%#v\n", err)
		}
		fmt.Printf("user:%#v\n", &user)
	}

	{
		var user User
		db.Model(&User{}).Preload("Courses").Find(&user, "id = ?", 24)
		fmt.Printf("user:%#v\n", &user)
		var courses []Course
		db.Model(&Course{}).Find(&courses)
		for _, course := range courses {
			fmt.Printf("course:%#v\n", &course)
		}
		err := db.Model(&user).Association("Courses").Replace(&courses[0], &courses[2])
		if err != nil {
			fmt.Printf("err:%#v\n", err)
		}
		fmt.Printf("user:%#v\n", &user)
	}
}

func TestAssociationDelete2(t *testing.T) {
	var user User
	db.Model(&User{}).Preload("Company").Preload("Courses").Find(&user, "id = ?", 25)
	err := db.Model(&user).Association("Company").Delete(&user.Company)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
	}
	err = db.Model(&user).Association("Courses").Delete(&user.Courses)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
	}
}

func TestAssociationClear(t *testing.T) {
	var user User
	db.Model(&User{}).Find(&user, "id = ?", 26)
	err := db.Model(&user).Association("Courses").Clear()
	if err != nil {
		fmt.Printf("err:%#v\n", err)
	}
}

func TestAssociationCount(t *testing.T) {
	var user User
	db.Model(&User{}).Find(&user, "id = ?", 27)
	count := db.Model(&user).Association("Company").Count()
	fmt.Printf("count:%#v\n", count)
	count = db.Model(&user).Association("Courses").Count()
	fmt.Printf("count:%#v\n", count)
	count = db.Model(&user).Association("IdentityCard").Count()
	fmt.Printf("count:%#v\n", count)
}

func TestAssociationDeleteUnscoped(t *testing.T) {
	var user User
	db.Model(&User{}).Find(&user, "id = ?", 28)
	err := db.Model(&user).Unscoped().Association("Courses").Clear()
	if err != nil {
		fmt.Printf("err:%#v\n", err)
	}
}

func TestEmbeddedPreload(t *testing.T) {

}
