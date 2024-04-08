package crud

import (
	"database/sql"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestSaveMethod(t *testing.T) {
	var user = User{
		Model: gorm.Model{
			ID: 16,
		},
		Age: 30,
		ActivateAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	// Save 保存所有字段，即使未修改或者是空值，如果传入结构体，不包含主键ID，则会创建新纪录
	printResult(db.Save(&user))
	// Select 选择指定字段，区分大小写，传入字段名或者表列名
	printResult(db.Select("age", "activate_at").Save(&user))
	// 传入map，必须通过 Model|Table 指定数据库表，默认必须包含 Where 子句
	printResult(db.Model(User{}).Save(map[string]interface{}{
		"age": 30,
		"ActivateAt": sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}))
	// 通过 session 配置开启全局批量更新
	// UPDATE `t_user` SET `activate_at`='2024-04-03 09:58:02.214',`age`=30,`updated_at`='2024-04-03 09:58:02.215' WHERE `t_user`.`deleted_at` IS NULL
	printResult(db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(User{}).Save(map[string]interface{}{
		"age": 30,
		"ActivateAt": sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}))
}

func TestUpdateMethod(t *testing.T) {
	// Update 更新单列
	// 默认禁用全局批量更新，必须搭配 Where 子句使用
	printResult(db.Model(&User{}).Update("age", 50))
	// 通过 session 配置开启全局批量更新
	printResult(db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(User{}).Update("age", 50))
	printResult(db.Model(&User{}).Where("id", 10).Update("age", 50))
}

func TestUpdatesMethod(t *testing.T) {
	var user = User{
		Model: gorm.Model{
			ID: 16,
		},
		Age: 30,
		ActivateAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	// Updates 更新多列，只更新非零值
	// UPDATE `t_user` SET `updated_at`='2024-04-03 10:04:11.484',`age`=30,`activate_at`='2024-04-03 10:04:11.483' WHERE `t_user`.`deleted_at` IS NULL AND `id` = 16
	printResult(db.Updates(&user))
	// 传入结构体必须包含主键信息，或者包含 Where 子句
	printResult(db.Updates(&User{Age: 80}))
	// 传入map 必须包含 Where 子句
	printResult(db.Model(&User{}).Updates(map[string]interface{}{
		"age": 16,
	}))
	// 通过 session 配置开启全局批量更新
	printResult(db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(User{}).Updates(map[string]interface{}{
		"age": 16,
	}))
	// Select 选择更新指定列
	// UPDATE `t_user` SET `age`=16,`updated_at`='2024-04-03 10:10:56.489' WHERE `id` = 16 AND `t_user`.`deleted_at` IS NULL
	printResult(db.Model(&User{}).Select("age").Where("id", 16).Updates(map[string]interface{}{
		"age":  16,
		"name": "rabbit",
	}))
}

func TestUpdateColumnMethod(t *testing.T) {
	// 跳过 Hook 函数 以及 时间追踪（updated_at）
	printResult(db.Model(&User{}).Where("id", 16).UpdateColumn("age", 19))
}

func TestUpdatingForCheckValueChanged(t *testing.T) {
	// 通过 hook 函数检测
	printResult(db.Model(&User{}).Where("id", 16).Updates(map[string]interface{}{
		"age":  46,
		"name": "rabbit",
	}))
}

func TestUpdatingForChangeValue(t *testing.T) {
	// 通过 hook 函数检测
	printResult(db.Model(&User{}).Where("id", 16).Updates(map[string]interface{}{
		"age": 46,
	}))
}
