package db

import (
	"fmt"
	"time"

	"blog/internal/config"
	"blog/internal/logger"

	"github.com/rabbit-rm/rabbit/componment/log/zapKit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func MySQL() *gorm.DB {
	return db
}

func init() {
	var err error
	username := config.MySQLConf().Username
	password := config.MySQLConf().Password
	host := config.MySQLConf().Host
	schemaName := config.MySQLConf().Schema
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, schemaName)
	if db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
		Logger: gormLogger.New(zapKit.NewZapMySQL(logger.L()), gormLogger.Config{
			SlowThreshold: 200 * time.Millisecond,
			Colorful:      true,
			LogLevel:      gormLogger.Error,
		}),
	}); err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute)

	/*db.Set("gorm:table_options", "ENGINE=InnoDB")
	if err := AutoMigrator(db.Migrator()); err != nil {
		panic(err)
	}*/

}

func AutoMigrator(migrator gorm.Migrator) error {
	/*if migrator.HasTable(models.Auth{}) {
		if err := migrator.DropTable(models.Auth{}); err != nil {
			return err
		}
	}
	if err := migrator.CreateTable(models.Auth{}); err != nil {
		return err
	}
	if migrator.HasTable(models.Tag{}) {
		if err := migrator.DropTable(models.Tag{}); err != nil {
			return err
		}
	}
	if err := migrator.CreateTable(models.Tag{}); err != nil {
		return err
	}
	if migrator.HasTable(models.Article{}) {
		if err := migrator.DropTable(models.Article{}); err != nil {
			return err
		}
	}
	if err := migrator.CreateTable(models.Article{}); err != nil {
		return err
	}
	return migrator.AutoMigrate(models.Tag{}, models.Article{}, models.Auth{})*/
	return nil
}
