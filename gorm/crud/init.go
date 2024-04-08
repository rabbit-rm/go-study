package crud

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
	port := 3306
	schema := "test"
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, schema)
	// logger
	var l logger.Interface
	// 默认日志
	// l = logger.Default.LogMode(logger.Info)
	// 慢查询打印
	l = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 3 * time.Millisecond,
		LogLevel:      logger.Info,
	})
	var err error
	if db, err = gorm.Open(mysql.Open(dns), &gorm.Config{Logger: l}); err != nil {
		log.Fatalf("init db error:%+v\n", err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("init sql db error:%+v\n", err)
	}
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetMaxIdleConns(100)
	sqlDb.SetConnMaxLifetime(10 * time.Second)
}

func setConnect(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(100)
	// 设置空闲连接最大存活时间
	sqlDB.SetConnMaxLifetime(10 * time.Second)
	return nil
}
