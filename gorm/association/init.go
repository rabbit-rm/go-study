package association

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	user := "root"
	password := "root@123456"
	host := "192.168.204.129"
	schema := "gorm_association_test"
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, 3306, schema)
	var err error
	var l logger.Interface
	l = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: time.Second,
		LogLevel:      logger.Info,
	})
	if db, err = gorm.Open(mysql.New(mysql.Config{DSN: dns}), &gorm.Config{Logger: l}); err != nil {
		log.Fatalf("connect mysql database error:%+v\n", err)
	}
	sdb, err := db.DB()
	if err != nil {
		log.Fatalf("database sql error:%+v\n", err)
	}
	sdb.SetMaxIdleConns(100)
	sdb.SetMaxOpenConns(100)
	sdb.SetConnMaxLifetime(10 * time.Second)
}

func printResult(tx *gorm.DB) {
	fmt.Printf("RowsAffected:%d,Error:%+v\n", tx.RowsAffected, tx.Error)
}
