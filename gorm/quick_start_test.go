package gormT

import (
	"fmt"
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TestStructModel struct {
	gorm.Model
	Field1              string `gorm:"<-:create;NOT NULL"`          // 允许读和创建
	Field2              string `gorm:"<-:update;NOT NULL"`          // 允许读和更新
	Field3              string `gorm:"<-;NOT NULL"`                 // 允许读写
	Field4              string `gorm:"<-:false;NOT NULL"`           // 允许读，禁止写
	Field5              string `gorm:"->;NOT NULL"`                 // 只读
	Field6              string `gorm:"->;<-:create;NOT NULL"`       // 允许读和创建
	Field7              string `gorm:"->:false;<-:create;NOT NULL"` // 只允许创建
	Field8              int64  `gorm:"column:age;check:age > 10;NOT NULL"`
	Field9              uint   `gorm:"autoCreateTime:milli;NOT NULL"`
	Field10             string `gorm:"-;NOT NULL"`           // 通过 struct 读写会忽略该字段
	Field11             string `gorm:"-:all;NOT NULL"`       // 通过 struct 读写 迁移会忽略该字段
	Field12             string `gorm:"-:migration;NOT NULL"` // 通过 struct 迁移会忽略该字段
	EmbeddedStructModel `gorm:"embedded;embeddedPrefix:embedded_"`
}

type EmbeddedStructModel struct {
	Field1 string `gorm:"column:A"`
	Field2 string `gorm:"column:B"`
	Field3 string `gorm:"column:C"`
	Field4 string `gorm:"column:D"`
}

func TestModelTag(t *testing.T) {
	dns := fmt.Sprintf("root:root@123456@tcp(192.168.204.129:3306)/test?charset=utf8mb4&parseTime=true&loc=Local")
	db, err := gorm.Open(mysql.Open(dns))
	if err != nil {
		log.Fatal(err)
	}
	if db.Migrator().HasTable(TestStructModel{}) {
		db.Migrator().DropTable(TestStructModel{})
	}
	db.Migrator().CreateTable(TestStructModel{})
}

func TestQuickStart(t *testing.T) {
	// 定义结构体映射表结构
	// Product 在 Gorm 中称之为 Model，一个 Model 对应一张数据库表，一个结构体实例对象对应一条数据库记录
	type Product struct {
		gorm.Model
		Code  string
		Price uint
	}

	// 连接数据库
	dns := fmt.Sprintf("root:root@123456@tcp(192.168.204.129:3306)/test?charset=utf8mb4&parseTime=true&loc=Local")
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: dns}))
	if err != nil {
		log.Fatal(err)
	}
	// 表存在则删除
	if db.Migrator().HasTable(&Product{}) {
		_ = db.Migrator().DropTable(&Product{})
	}
	// 自动迁移表结构
	db.AutoMigrate(&Product{})
	// 增加数据
	db.Create(&Product{
		Code:  "D42",
		Price: 100,
	})
	db.Create(&Product{
		Code:  "D45",
		Price: 110,
	})
	db.Create(&Product{
		Code:  "D52",
		Price: 150,
	})
	{
		// 查找数据
		{
			var p Product
			db.First(&p, 1)
			fmt.Printf("product:{ID:%d,Code:%s,Price:%d}\n", p.ID, p.Code, p.Price)
		}
		{
			var p Product
			db.First(&p, "code = ?", "D52")
			fmt.Printf("product:{ID:%d,Code:%s,Price:%d}\n", p.ID, p.Code, p.Price)
		}
	}
	{
		// 更新数据
		{
			var p Product
			db.First(&p, 1)
			db.Model(&p).Update("Price", 200)
			fmt.Printf("product:{ID:%d,Code:%s,Price:%d}\n", p.ID, p.Code, p.Price)
		}
		{
			var p Product
			db.First(&p, "code = ?", "D45")
			db.Model(&p).Updates(Product{
				Code:  "F45",
				Price: 210,
			})
			fmt.Printf("product:{ID:%d,Code:%s,Price:%d}\n", p.ID, p.Code, p.Price)
		}
		{
			var p Product
			db.First(&p, "code = ?", "D52")
			db.Model(&p).Updates(map[string]interface{}{
				"price": 350,
			})
			fmt.Printf("product:{ID:%d,Code:%s,Price:%d}\n", p.ID, p.Code, p.Price)
		}
	}
	{
		// 删除数据
		{
			var p Product
			// 逻辑删除
			db.Delete(&p, 1)
		}
	}

}
