package crud

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func mockUser(db *gorm.DB) {
	var users []User
	for i := range 10000 {
		users = append(users, User{
			Name:     fmt.Sprintf("user-%d", i),
			Email:    fmt.Sprintf("user-%d@gmail.com", i),
			Age:      uint8(20 + i%10),
			Birthday: time.Now(),
		})
	}
	tx := db.CreateInBatches(&users, 100)
	fmt.Printf("RowsAffected:%d\n", tx.RowsAffected)
	fmt.Printf("Error:%v\n", tx.Error)
}
