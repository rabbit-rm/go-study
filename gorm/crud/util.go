package crud

import (
	"fmt"

	"gorm.io/gorm"
)

func printResult(tx *gorm.DB) {
	fmt.Printf("RowsAffected:%d,Error:%+v\n", tx.RowsAffected, tx.Error)
}
