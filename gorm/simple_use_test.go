package gormT

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

func TestCreateSingleRecord(t *testing.T) {
	user := &User{
		Name:     "rabbit",
		Email:    "rabbit.rm99@gmail.com",
		Age:      24,
		Birthday: time.Now(),
	}
	{
		// 插入单条记录
		tx := db.Create(user)
		fmt.Printf("user:%+v\n", user)
		fmt.Printf("rows affected:%d\n", tx.RowsAffected)
		fmt.Printf("error:%+v\n", tx.Error)
	}
}

func TestCreateMultiRecord(t *testing.T) {
	users := []*User{
		{
			Name:     "rabbit1",
			Email:    "rabbit1@gmail.com",
			Age:      24,
			Birthday: time.Now(),
		},
		{
			Name:     "rabbit2",
			Email:    "rabbit2@gmail.com",
			Age:      24,
			Birthday: time.Now(),
		},
		{
			Name:     "rabbit3",
			Email:    "rabbit3@gmail.com",
			Age:      24,
			Birthday: time.Now(),
		},
	}
	{
		// 插入多条记录
		tx := db.Create(users)
		fmt.Printf("user:%+v\n", users)
		fmt.Printf("rows affected:%d\n", tx.RowsAffected)
		fmt.Printf("error:%+v\n", tx.Error)
	}
}

func TestCreateInBatches(t *testing.T) {
	type Product struct {
		gorm.Model
		Code  string
		Price uint
	}
	p := Product{
		Code:  "D10N",
		Price: 100,
	}
	pp := []Product{p, p, p, p, p}
	tx := db.CreateInBatches(pp, 2)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestCreateSQLExp(t *testing.T) {
	tx := db.Table("products").Create(map[string]interface{}{
		"code": "D20",
		// 等价于 gorm.Expr
		"price": clause.Expr{SQL: "GREATEST(?,?,?,?,?)", Vars: []interface{}{1, 2, 3, 4, 56}},
	})
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)

	db.Clauses()
}

func TestSelectFirst(t *testing.T) {
	var u User
	tx := db.First(&u)
	fmt.Printf("user:%+v\n", u)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestSelectLast(t *testing.T) {
	var u User
	tx := db.Last(&u)
	fmt.Printf("user:%+v\n", u)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestSelectWhere(t *testing.T) {
	var u User
	tx := db.Where("name = ?", "rabbit2").Find(&u)
	fmt.Printf("user:%+v\n", u)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestSelectAnd(t *testing.T) {
	var users []map[string]interface{}
	// tx := db.Table("t_user").Select("name", "email", "age").Where("name LIKE ?", "rabbit%").Where("age > ?", 30).Scan(&users)
	tx := db.Table("t_user").Select("name", "email", "age").Where("name LIKE ? OR age > ?", "rabbit%", 30).Scan(&users)
	fmt.Printf("user:%+v\n", users)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestSelectOrder(t *testing.T) {
	var users []map[string]interface{}
	tx := db.Table("t_user").Select("name", "age").Order("age desc").Scan(&users)
	fmt.Printf("user:%+v\n", users)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestSelectWhereLike(t *testing.T) {
	var users []User
	tx := db.Model(User{}).Where("name LIKE ?", "rabbit%").Find(&users)
	fmt.Printf("user:%+v\n", users)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestSelectPart(t *testing.T) {
	var u User
	tx := db.Select("name", "age").Where("ID = ?", 1).Find(&u)
	fmt.Printf("user:%s\n", u.String())
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestOrder(t *testing.T) {
	var users []*User
	tx := db.Select("ID", "name", "age", "email").Order("ID desc").Find(&users)
	fmt.Printf("users:%v\n", users)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestLimit(t *testing.T) {
	var users []*User
	tx := db.Select("ID", "name", "age", "email").Order("ID desc").Limit(2).Find(&users)
	fmt.Printf("users:%v\n", users)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}
func TestOffset(t *testing.T) {
	var users []*User
	tx := db.Select("ID", "name", "age", "email").Order("ID desc").Limit(1).Offset(0).Find(&users)
	fmt.Printf("users:%v\n", users)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestLimitOffset(t *testing.T) {
	var users []*User
	var users2 []*User
	// Limit -1 取消 limit 限制
	// offset -1 取消 offset 限制
	tx := db.Select("ID", "name", "age", "email").Order("ID desc").Limit(1).Offset(1).Find(&users).Limit(-1).Offset(-1).Find(&users2)
	fmt.Printf("users:%v\n", users)
	fmt.Printf("users2:%v\n", users2)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestCount(t *testing.T) {
	var c int64
	tx := db.Model(&User{}).Count(&c)
	fmt.Printf("count:%d\n", c)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestNestedSubQuery(t *testing.T) {
	var avgs []float64
	sub := db.Select("AVG(age)").Where("name LIKE ?", "rabbit%").Model(&User{})
	// having 子查询
	tx := db.Model(&User{}).Select("AVG(age) as averages").Group("name").Having("AVG(age) >= (?)", sub).Find(&avgs)
	fmt.Printf("avgs:%v\n", avgs)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestAdvancedSelect(t *testing.T) {
	type UserAPI struct {
		ID   uint
		Name string
	}
	var user UserAPI
	// tx := db.Model(&User{}).First(&user)
	tx := db.Session(&gorm.Session{QueryFields: false}).Model(User{}).First(&user)
	fmt.Printf("user:%+v\n", user)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestSave(t *testing.T) {
	var user User
	tx := db.Where("ID", 1).Find(&user)
	fmt.Printf("user:%v\n", &user)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
	user.Age = 25
	// 保存对应 Model,如果主键ID存在就 Update，否则就 insert
	tx = db.Save(&user)
	fmt.Printf("user:%v\n", &user)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)

	var user2 = &User{
		Name:     "user-2",
		Email:    "user@gmail.com",
		Age:      34,
		Birthday: time.Now(),
	}
	// INSERT
	tx = db.Save(&user2)
	fmt.Printf("user2:%v\n", user2)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestUpdate(t *testing.T) {
	var user User
	tx := db.Where("name = ?", "rabbit").Find(&user)
	fmt.Printf("user:%v\n", &user)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
	tx = db.Model(user).Update("age", 25)
	fmt.Printf("user:%v\n", &user)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)

	// 更新多字段
	tx = db.Model(user).Updates(&User{Age: 25, Birthday: time.Now()})
	fmt.Printf("user:%v\n", &user)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestExpr(t *testing.T) {
	var user User
	tx := db.Where("name = ?", "rabbit1").Find(&user)
	fmt.Printf("user:%v\n", &user)
	// expr SQL 表达式
	db.Model(&user).Updates(map[string]interface{}{
		"age":      gorm.Expr("age + ?", 1),
		"birthday": (&user).Birthday.Add(time.Hour * 24),
	})
	fmt.Printf("user:%v\n", &user)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestDelete(t *testing.T) {
	var user User
	tx := db.Where("name = ?", "rabbit1").Find(&user)
	fmt.Printf("user:%v\n", &user)
	// gorm 采用逻辑删除，删除数据时 只更新 delete_at 字段，所有查询语句都会加上 WHERE delete_at IS NULL 条件
	tx = tx.Delete(&user)
	fmt.Printf("user:%v\n", &user)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestUnscoped(t *testing.T) {
	var user User
	// Unscoped 跳过逻辑作用域
	tx := db.Unscoped().Where("name = ?", "rabbit1").Find(&user)
	fmt.Printf("user:%v\n", &user)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)

	// Unscoped 跳过逻辑作用域，物理删除
	tx = db.Unscoped().Delete(&User{}, "name = ?", "rabbit1")
	fmt.Printf("user:%v\n", &user)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

// ========================================
//        多表查询
// ========================================

type Comment struct {
	gorm.Model
	Content string `gorm:"column:content"`
	PostID  uint   `gorm:"column:post_id"`
	Post    *Post
}

func (Comment) TableName() string {
	return "t_comment"
}

type Tag struct {
	gorm.Model
	Name  string  `gorm:"column:name"`
	Posts []*Post `gorm:"many2many:t_post_tags"`
}

func (Tag) TableName() string {
	return "t_tag"
}

type Post struct {
	gorm.Model
	Title    string     `gorm:"column:title"`
	Content  string     `gorm:"column:content"`
	Comments []*Comment `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;references:ID"`
	Tags     []*Tag     `gorm:"many2many:t_post_tags"`
}

func (Post) TableName() string {
	return "t_post"
}

func TestAssociationCreate(t *testing.T) {
	var p Post
	p = Post{
		Title:   "post1",
		Content: "content1",
		Comments: []*Comment{
			{
				Content: "comment1",
				Post:    &p,
			},
			{
				Content: "comment2",
				Post:    &p,
			},
		},
		Tags: []*Tag{
			{
				Name: "tag1",
			},
			{
				Name: "tag2",
			},
		},
	}
	tx := db.Create(&p)
	fmt.Printf("post:%v\n", &p)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestAssociationSave(t *testing.T) {
	var p Post
	p = Post{
		Title:   "post2",
		Content: "content2",
		Comments: []*Comment{
			{
				Content: "comment3",
				Post:    &p,
			},
			{
				Content: "comment4",
				Post:    &p,
			},
		},
		Tags: []*Tag{
			{
				Name: "tag1",
			},
			{
				Name: "tag3",
			},
		},
	}
	tx := db.Save(&p)
	fmt.Printf("post:%v\n", &p)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestAssociationSelect(t *testing.T) {
	var (
		p        Post
		comments []Comment
	)
	p.ID = 1
	// Post 是源模型，主键 ID 不能为空
	// Association 方法指定关联字段名，在 Post 模型中关联的评论使用 Comments 表示
	// 最后使用 Find 方法来查询关联的评论
	err := db.Model(&p).Association("Comments").Find(&comments)
	if err != nil {
		log.Fatalf("err:%+v\n", err)
	}
	fmt.Printf("post:%v\n", &p)
	fmt.Printf("comments:%v\n", comments)
}

func TestAssociationPreload(t *testing.T) {
	var (
		p Post
	)
	p.ID = 1
	// Preload 预加载关联字段名，这样在查询 Post 时，自动填充关联字段
	tx := db.Preload("Comments").Preload("Tags").Find(&p)
	fmt.Printf("post:%v\n", &p)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestAssociationSelectJoin(t *testing.T) {
	type PostComment struct {
		Title   string
		Comment string
	}
	var p = Post{}
	p.ID = 1
	pc := PostComment{}
	tx := db.Model(&p).Select("t_post.Title,t_comment.Content as comment").Joins("LEFT JOIN t_comment ON t_comment.post_id = t_post.id").Scan(&pc)
	fmt.Printf("post comment:%v\n", &pc)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestAssociationJoin(t *testing.T) {
	var cc []*Comment
	// SELECT `t_comment`.`id`,
	//       `t_comment`.`created_at`,
	//       `t_comment`.`updated_at`,
	//       `t_comment`.`deleted_at`,
	//       `t_comment`.`content`,
	//       `t_comment`.`post_id`,
	//       `t_post`.`id`         AS `Post__id`,
	//       `t_post`.`created_at` AS `Post__created_at`,
	//       `t_post`.`updated_at` AS `t_Post__updated_at`,
	//       `t_post`.`deleted_at` AS `Post__deleted_at`,
	//       `t_post`.`title`      AS `Post__title`,
	//       `t_post`.`content`    AS `Post__content`
	// FROM `t_comment`
	//         LEFT JOIN `t_post` ON `t_comment`.`post_id` = `t_post`.`id` AND `t_post`.`deleted_at` IS NULL
	// WHERE `t_comment`.`deleted_at` IS NULL
	tx := db.Joins("Post").Find(&cc)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
	for i, c := range cc {
		fmt.Printf("[%d]{content:%s,PostID:%d,{title:%s,content:%s}}\n", i, c.Content, c.PostID, c.Post.Title, c.Post.Content)
	}

}

func TestAssociationUpdate(t *testing.T) {
	// var p Post
	// tx := db.First(&p)
	// fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	// fmt.Printf("error:%+v\n", tx.Error)
	// var comments []*Comment
	// tx = db.Where("post_id = ?", p.ID).Find(&comments)
	// fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	// fmt.Printf("error:%+v\n", tx.Error)
	// for i, comment := range comments {
	// 	comment.Content = fmt.Sprintf("new comment content %d", i)
	// 	tx = db.Model(&Comment{}).Where(comment.ID).Updates(&comment)
	// 	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	// 	fmt.Printf("error:%+v\n", tx.Error)
	// }

	var p Post
	tx := db.First(&p)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
	var comments = []*Comment{
		{
			Content: "New Comment Connect 1",
		},
		{
			Content: "New Comment Connect 2",
		},
	}
	// BEGIN TRANSACTION;
	// INSERT INTO `comment` (`created_at`,`updated_at`,`deleted_at`,`content`,`post_id`) VALUES ('2023-05-23 09:07:42.852','2023-05-23 09:07:42.852',NULL,'comment3',1) ON DUPLICATE KEY UPDATE `post_id`=VALUES(`post_id`)
	// UPDATE `post` SET `updated_at`='2023-05-23 09:07:42.846' WHERE `post`.`deleted_at` IS NULL AND `id` = 1
	// UPDATE `comment` SET `post_id`=NULL WHERE `comment`.`id` <> 8 AND `comment`.`post_id` = 1 AND `comment`.`deleted_at` IS NULL
	// COMMIT;
	// 以下代码执行上面 SQL，replace 替换 postId == 1 的所有 comment(新建记录)，取消原有 comment 的外键关系，不会删除原 comment
	err := db.Model(&p).Association("Comments").Replace(comments)
	if err != nil {
		log.Fatal(err)
	}
}

func TestAssociationDelete(t *testing.T) {
	var p Post
	tx := db.Where("id = ?", 1).Delete(&p)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)

	// association delete

	var pp Post
	tx = db.Where("id = ?", 3).Find(&p)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
	err := db.Model(&pp).Association("Comments").Delete(&pp.Comments)
	if err != nil {
		log.Fatal(err)
	}
}

func TestJoins(t *testing.T) {
	var res []map[string]interface{}
	tx := db.Table("t_post p").Select("p.id pid,p.title ptitle,p.content pcontent,c.id cid,c.content ccontent").Joins("LEFT JOIN t_comment c on p.id = c.post_id").Scan(&res)
	fmt.Printf("res:%+v\n", res)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
	rows, err := db.Table("t_post p").Joins("LEFT JOIN t_comment c on p.id = c.post_id").Rows()
	if err != nil {
		log.Fatalf("err:%+v\n", err)
	}
	for rows.Next() {
		var a interface{}
		var b interface{}
		var c interface{}
		var d interface{}
		var e interface{}
		var f interface{}
		var g interface{}
		var h interface{}
		var i interface{}
		var j interface{}
		var k interface{}
		var l interface{}
		err := rows.Scan(&a, &b, &c, &d, &e, &f, &g, &h, &i, &j, &k, &l)
		if err != nil {
			log.Fatalf("err:%+v\n", err)
		}
		fmt.Printf("a:%v,b:%s,c:%s,d:%v,e:%v,f:%v,g:%v,h:%s,i:%v,j:%v,k:%v,l:%v\n", a, b, c, d, e, f, g, h, i, j, k, l)
	}
}

// ========================================
//
//	事务
//
// ========================================
func TestTransactionPost(t *testing.T) {
	// Transaction 自动管理事务
	// 接收一个函数，真个函数内部都将在一个事务中被处理，返回 nil 事务提交，返回 err 事务回滚
	err := db.Transaction(func(tx *gorm.DB) error {
		p := Post{
			Title: "Transaction Post 1",
		}
		if err := tx.Create(&p).Error; err != nil {
			return err
		}
		comment := Comment{Content: "Transaction Comment Connect 1", PostID: p.ID}
		if err := tx.Create(&comment).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestTransactionPostWithManually(t *testing.T) {
	// 手动管理事务
	// 开始事务，返回 tx 后续都应该使用 tx
	tx := db.Begin()
	p := Post{
		Title: "Manually Transaction Post",
	}
	if err := tx.Create(&p).Error; err != nil {
		// 事务回滚
		tx.Rollback()
		log.Fatal(err)
	}
	comment := Comment{Content: "Manually Transaction Comment Connect", PostID: p.ID}
	if err := tx.Create(&comment).Error; err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	// 事务提交
	if tx.Commit().Error != nil {
		tx.Rollback()
	}
}

func TestHook(t *testing.T) {
	u := User{
		Name:     "user-10",
		Email:    "user-10@gamil.com",
		Age:      7,
		Birthday: time.Now(),
	}
	tx := db.Create(&u)
	fmt.Printf("user:%s\n", (&u).String())
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestRow(t *testing.T) {
	var res = make(map[string]interface{})
	// 原生查询 sql
	tx := db.Raw("SELECT id,name,age FROM t_user WHERE id = ?", 1).Scan(res)
	fmt.Println(res)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestRowExpr(t *testing.T) {
	var res = make(map[string]interface{})
	tx := db.Raw("SELECT SUM(age) as sumage FROM t_user WHERE member_number ?", gorm.Expr("IS NULL")).Scan(res)
	fmt.Println(res)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestExec(t *testing.T) {
	db.Exec("UPDATE t_user SET age = ? WHERE id IN ?", 18, []int64{1, 2, 3, 4, 5})
	db.Exec("UPDATE t_user SET age = ? WHERE id IN ?", gorm.Expr("age * ?", 2), []int64{6, 7, 8, 9, 10})
}

func TestExecNamed(t *testing.T) {
	var res = make([]map[string]interface{}, 0)
	tx := db.Raw("SELECT name,email FROM t_user WHERE name LIKE @name", sql.Named("name", "rabbit%")).Scan(&res)
	fmt.Println(res)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestDryRun(t *testing.T) {
	var u User
	// dryRun 空跑 不执行sql，可用于调试生成的 sql
	tx := db.Session(&gorm.Session{DryRun: true}).First(&u, 1)
	fmt.Printf("sql：%v\n", tx.Statement.SQL.String())
	fmt.Printf("sql vars：%v\n", tx.Statement.Vars)

	tx = db.Session(&gorm.Session{DryRun: true}).Exec("UPDATE t_user SET age = ? WHERE id IN ?", gorm.Expr("age * ?", 2), []int64{6, 7, 8, 9, 10})
	fmt.Printf("sql：%v\n", tx.Statement.SQL.String())
	fmt.Printf("sql vars：%v\n", tx.Statement.Vars)
}

func TestNotFound(t *testing.T) {
	/*var user User
	tx := db.Model(User{}).FirstOrInit(&user)
	fmt.Printf("user:%v\n", user)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)*/
	db.Assign()
	var user map[string]interface{} = make(map[string]interface{})
	tx := db.Table("users").FirstOrInit(&user, map[string]interface{}{
		"name": "fff",
	})
	fmt.Printf("user:%v\n", user)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestHint(t *testing.T) {
	var u User
	_ = db.Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find(&u)

}

func TestFindInBatches(t *testing.T) {
	var users []User
	db.Model(User{}).FindInBatches(&users, 10, func(tx *gorm.DB, batch int) error {
		fmt.Printf("users:%d,batch:%d\n", len(users), batch)
		return nil
	})
}

func TestPluck(t *testing.T) {
	var names []string
	tx := db.Model(User{}).Where("id > 1000 AND id < 2000").Pluck("name", &names)
	fmt.Printf("names:%d\n", len(names))
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}

func TestScopes(t *testing.T) {
	// 定义 scopes
	idGt1000 := func(db *gorm.DB) *gorm.DB {
		return db.Where("id > 1000")
	}
	idLt2000 := func(db *gorm.DB) *gorm.DB {
		return db.Where("id < 2000")
	}
	var users []User
	tx := db.Model(User{}).Scopes(idLt2000, idGt1000).Find(&users)
	fmt.Printf("rows affected:%d\n", tx.RowsAffected)
	fmt.Printf("error:%+v\n", tx.Error)
}
