package association

import (
	"fmt"

	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name string `gorm:"COLUMN:name"`
}

func (*Company) TableName() string {
	return "t_company"
}

func (company *Company) Format(state fmt.State, verb rune) {
	_, _ = fmt.Fprintf(state, "{ID:%2d,Name:%s}", company.ID, company.Name)
}

type CompanyUser struct {
	Company
	Users []User
}

type IdentityCard struct {
	gorm.Model
	No      string `gorm:"COLUMN:no"`
	Address string `gorm:"COLUMN:address"`
}

func (*IdentityCard) TableName() string {
	return "t_identity_card"
}

func (identityCard *IdentityCard) Format(state fmt.State, verb rune) {
	_, _ = fmt.Fprintf(state, "{ID:%2d,No:%s,Address:%s}", identityCard.ID, identityCard.No, identityCard.Address)
}

type IdentityCardUser struct {
	IdentityCard
	User User
}

type CourseUser struct {
	Course
	Users []User `gorm:"many2many:t_course_user;foreignKey:ID;references:ID"`
}

type Course struct {
	gorm.Model
	Name string `gorm:"COLUMN:name"`
}

func (*Course) TableName() string {
	return "t_course"
}

func (course *Course) Format(state fmt.State, verb rune) {
	_, _ = fmt.Fprintf(state, "{ID:%2d,Name:%s}", course.ID, course.Name)
	return
}

type User struct {
	gorm.Model
	Name           string       `gorm:"COLUMN:name"`
	CompanyID      uint         `gorm:"COLUMN:company_id"`
	Company        Company      `gorm:"foreignKey:CompanyID;references:ID"`
	IdentityCardID uint         `gorm:"COLUMN:identity_card_id"`
	IdentityCard   IdentityCard `gorm:"foreignKey:ID;references:ID"`
	Courses        []*Course    `gorm:"many2many:t_user_course;foreignKey:ID;references:ID"`
}

func (*User) TableName() string {
	return "t_user"
}

func (user *User) Format(state fmt.State, verb rune) {
	_, _ = fmt.Fprintf(state, "{ID:%2d,Name:%s,Company:%#v,IdentidyCard:%#v,courses:%#v}", user.ID, user.Name, &user.Company, &user.IdentityCard, user.Courses)
	return
}

/*type CreditCard struct {
	gorm.Model
	No    uint    `gorm:"COLUMN:no"`
	Quota float64 `gorm:"COLUMN:quota"`
}

type CreditCardUser struct {
	CreditCard
	User User
}

func (*CreditCard) TableName() string {
	return "t_credit_card"
}*/
