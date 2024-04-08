package crud

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rabbit-rm/rabbit/errorKit"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string         `gorm:"column:name;NOT NULL;comment:用户名"`
	Email        string         `gorm:"column:email;NOT NULL;unique;comment:邮箱"`
	Age          uint8          `gorm:"column:age;NOT NULL;default:0;comment:年龄"`
	Birthday     time.Time      `gorm:"column:birthday;NOT NULL;comment:生日"`
	MemberNumber sql.NullString `gorm:"column:member_number;comment:成员编号"`
	ActivateAt   sql.NullTime   `gorm:"column:activate_at;comment:激活时间"`
}

func (user *User) TableName() string {
	return "t_user"
}

func (user *User) String() string {
	timeFormat := func(t time.Time) string {
		return t.Format("2006-01-02T15:04:05.000")
	}
	buffer := bytes.NewBuffer(make([]byte, 0))
	buffer.WriteString("{")
	if user.ID != 0 {
		// write ID
		buffer.WriteString(fmt.Sprintf("ID:%d,", user.ID))
	}
	if strings.Compare(user.Name, "") != 0 {
		// write name
		buffer.WriteString(fmt.Sprintf("name:%s,", user.Name))
	}
	if strings.Compare(user.Email, "") != 0 {
		// write name
		buffer.WriteString(fmt.Sprintf("email:%s,", user.Email))
	}
	if user.Age != 0 {
		// write age
		buffer.WriteString(fmt.Sprintf("age:%d,", user.Age))
	}
	if user.Birthday.Sub(time.Time{}) != 0 {
		// write birthday
		buffer.WriteString(fmt.Sprintf("birthday:%s,", timeFormat(user.Birthday)))
	}
	if user.MemberNumber.Valid {
		// write MemberNumber
		buffer.WriteString(fmt.Sprintf("memberNumber:%s,", user.MemberNumber.String))
	}
	if user.ActivateAt.Valid {
		// write ActivateAt
		buffer.WriteString(fmt.Sprintf("activateAt:%s,", timeFormat(user.ActivateAt.Time)))
	}
	if user.CreatedAt.Sub(time.Time{}) != 0 {
		// write CreatedAt
		buffer.WriteString(fmt.Sprintf("createdAt:%s,", timeFormat(user.CreatedAt)))
	}
	if user.UpdatedAt.Sub(time.Time{}) != 0 {
		// write UpdatedAt
		buffer.WriteString(fmt.Sprintf("updatedAt:%s,", timeFormat(user.UpdatedAt)))
	}
	if user.DeletedAt.Valid {
		// write DeletedAt
		buffer.WriteString(fmt.Sprintf("deletedAt:%s,", timeFormat(user.DeletedAt.Time)))
	}
	buffer.Truncate(buffer.Len() - 1)
	buffer.WriteString("}")
	return buffer.String()
}

func (user *User) Format(state fmt.State, verb rune) {
	_, _ = io.WriteString(state, user.String())
}

// BeforeSave	调用 Save 前
// AfterSave	调用 Save 后
// AfterCreate	插入记录后
// BeforeUpdate	更新记录前
// AfterUpdate	更新记录后
// BeforeDelete	删除记录前
// AfterDelete	删除记录后

// BeforeCreate	插入记录前
func (user *User) BeforeCreate(tx *gorm.DB) error {
	fmt.Printf("bofore create check age > 8\n")
	if user.Age <= 8 {
		return errorKit.New("age < 8")
	}
	tx.Statement.Changed()
	return nil
}

// AfterFind	查询记录后
func (user *User) AfterFind(tx *gorm.DB) error {
	fmt.Printf("user:%s\n", user)
	return nil
}

func (user *User) BeforeUpdate(tx *gorm.DB) error {
	// 不允许修改名称
	if tx.Statement.Changed("name") {
		return errorKit.New("name not allow to change")
	}
	if tx.Statement.Changed("age") {
		tx.Statement.SetColumn("age", 80)
	}
	return nil
}

type ProductInventory struct {
	gorm.Model
	ProductID uint64        `gorm:"column:product_id"`
	Quantity  sql.NullInt64 `gorm:"column:quantity"`
}

func (p ProductInventory) TableName() string {
	return "t_product_inventory"
}

func createUserTable(db *gorm.DB, model interface{}) error {
	migrator := db.Migrator()
	if !migrator.HasTable(model) {
		if err := migrator.CreateTable(model); err != nil {
			return err
		}
	}
	return nil
}
