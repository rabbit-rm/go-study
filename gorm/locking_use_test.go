package gormT

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestLockingUpdate(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(db *gorm.DB, index int) {
			defer wg.Done()
			err := db.Transaction(func(tx *gorm.DB) error {
				var p ProductInventory
				// 锁定当前行，阻止其他事务修改数据
				tx = tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).Where("product_id", 1).Find(&p)
				fmt.Printf("%s bofore update:{product_id:%d,quantity:%d}\n", fmt.Sprintf("[goroutine %d]", index), p.ProductID, p.Quantity.Int64)
				if index%2 == 0 {
					p.Quantity = sql.NullInt64{
						Int64: int64(index),
						Valid: true,
					}
					time.Sleep(time.Duration(index) * time.Second)
					tx = tx.Updates(&p)
				}
				// tx = tx.Where("product_id", 1).Find(&p)
				fmt.Printf("%s after update:{product_id:%d,quantity:%d}\n", fmt.Sprintf("[goroutine %d]", index), p.ProductID, p.Quantity.Int64)
				return nil
			})
			if err != nil {
				log.Fatalf("error:%+v\n", err)
			}
		}(db, i+1)
	}
	wg.Wait()
}

func TestLockingShare(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(db *gorm.DB) {
		defer wg.Done()
		err := db.Transaction(func(tx *gorm.DB) error {
			var p ProductInventory
			for {
				tx = tx.Where("product_id", 2).Find(&p)
				fmt.Printf("%s bofore update:{product_id:%d,quantity:%d}\n", fmt.Sprintf("[read]"), p.ProductID, p.Quantity.Int64)
				time.Sleep(time.Second)
			}
		})
		if err != nil {
			log.Fatalf("error:%+v\n", err)
		}
	}(db)
	wg.Add(1)
	go func(db *gorm.DB) {
		defer wg.Done()
		err := db.Transaction(func(tx *gorm.DB) error {
			var p ProductInventory
			time.Sleep(3 * time.Second)
			tx = tx.Where("product_id", 2).Delete(&p)
			fmt.Printf("%s bofore update:{product_id:%d,quantity:%d}\n", fmt.Sprintf("[delete]"), p.ProductID, p.Quantity.Int64)
			return nil
		})
		if err != nil {
			log.Fatalf("error:%+v\n", err)
		}
	}(db)
	wg.Add(1)
	go func(db *gorm.DB) {
		defer wg.Done()
		err := db.Transaction(func(tx *gorm.DB) error {
			var p ProductInventory
			// 锁定当前行，阻止其他事务修改数据
			tx = tx.Clauses(clause.Locking{Strength: clause.LockingStrengthShare}).Where("product_id", 2).Find(&p)
			fmt.Printf("%s bofore update:{product_id:%d,quantity:%d}\n", fmt.Sprintf("[update]"), p.ProductID, p.Quantity.Int64)
			p.Quantity = sql.NullInt64{
				Int64: int64(2),
				Valid: true,
			}
			time.Sleep(time.Duration(10) * time.Second)
			tx = tx.Updates(&p)
			fmt.Printf("%s after update:{product_id:%d,quantity:%d}\n", fmt.Sprintf("[update]"), p.ProductID, p.Quantity.Int64)
			return nil
		})
		if err != nil {
			log.Fatalf("error:%+v\n", err)
		}
	}(db)

	wg.Wait()
}
