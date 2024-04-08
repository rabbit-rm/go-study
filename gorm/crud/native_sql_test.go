package crud

import (
	"fmt"
	"testing"

	"gorm.io/gorm"
)

func TestNativeSQLRaw(t *testing.T) {
	var user User
	tx := db.Raw(`SELECT * FROM t_user WHERE id IN (?,?,?)`, 6, 7, 8)
	fmt.Println(tx.RowsAffected)
	tx.Scan(&user)
	fmt.Println(user)
}

func TestNativeSQLExec(t *testing.T) {
	var user User
	db.Association()
	tx := db.Model(&User{}).Find(&user, "id = ?", 20)
	tx = db.Exec(`DELETE FROM t_user WHERE id = ?`, 20)
	fmt.Println(tx.RowsAffected)
	fmt.Println(user)
}

func TestNativeDryRunMode(t *testing.T) {
	var user User
	db.Session(&gorm.Session{DryRun: true}).Delete(&user, "id = ?", 1)
}

func TestToSQL(t *testing.T) {
	db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		var user User
		return tx.Delete(&user, "id = ?", 1)
	})
}

func TestRows(t *testing.T) {
	db.Select(nil).Row()
	db.Select(nil).Rows()
}
