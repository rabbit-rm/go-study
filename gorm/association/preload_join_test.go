package association

import (
	"fmt"
	"testing"
)

func TestPreload(t *testing.T) {
	var user User
	db.Model(&User{}).Preload("Company").Preload("Courses").Preload("IdentityCard").Find(&user, "id = ?", 31)
	fmt.Printf("user:%#v\n", &user)
}

func TestJoins(t *testing.T) {
	var user User
	db.Joins("Company").Find(&user, "t_user.id = ?", 32)
	fmt.Printf("user:%#v\n", &user)
	db.Joins("IdentityCard").Find(&user, "t_user.id = ?", 32)
	fmt.Printf("user:%#v\n", &user)
}
