package association

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestError(t *testing.T) {
	var user User
	tx := db.Model(&User{}).Find(&user, "id = ")
	if err := tx.Error; err != nil {
		fmt.Printf("error: %#v\n", tx.Error)
		var merr *mysql.MySQLError
		if errors.As(err, &merr) {
			fmt.Printf("mysql error:%#v\n", merr.Error())
		}
	}
}
