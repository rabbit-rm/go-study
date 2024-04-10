package association

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestMockData(t *testing.T) {
	var user = User{
		Name: "rabbit",
		Company: Company{
			Name: "YozoSoft",
		},
		IdentityCard: IdentityCard{
			No:      "342401199999996337",
			Address: "JiangSu",
		},
		Courses: []*Course{
			{
				Name: "Java",
			},
			{
				Name: "Golang",
			},
		},
	}
	tx := db.Model(&User{}).Create(&user)
	printResult(tx)
}

func TestMockUserData(t *testing.T) {
	tx := db.Begin()
	for i := 3; i < 50; i++ {
		var user = User{
			Name: fmt.Sprintf("user%02d", i),
			Company: Company{
				Model: gorm.Model{
					ID: uint(i%4) + 1,
				},
			},
			IdentityCard: IdentityCard{
				No:      randomIdentityNo(),
				Address: randomIdentityAddress(),
			},
			Courses: []*Course{
				{
					Model: gorm.Model{
						ID: uint(i%4) + 1,
					},
				},
				{
					Model: gorm.Model{
						ID: uint(i%3) + 5,
					},
				},
			},
		}
		tx = tx.Model(&User{}).Create(&user)
		printResult(tx)
		if tx.Error != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
}

func TestMockCourseData(t *testing.T) {
	var courses = []Course{
		{
			Name: "C++",
		},
		{
			Name: "PHP",
		},
		{
			Name: "javascript",
		},
		{
			Name: "Vue",
		},
		{
			Name: "Ruby",
		},
	}
	printResult(db.Model(&Course{}).Create(&courses))
}

func TestMockCompanyData(t *testing.T) {
	var companyes = []Company{
		{
			Name: "JD",
		},
		{
			Name: "Alibaba",
		},
		{
			Name: "ByteDance",
		},
	}
	printResult(db.Model(&Company{}).Create(&companyes))
}

func randomIdentityNo() string {
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	limit := 18
	var no []byte
	for i := 0; i < limit; i++ {
		no = strconv.AppendInt(no, r.Int63n(10), 10)
	}
	return string(no)
}

func randomIdentityAddress() string {
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	var addresses = []string{
		"JiangSu",
		"Anhui",
		"ShangHai",
		"NanJ",
	}
	return addresses[r.Int63n(4)]
}

func TestName(t *testing.T) {
	fmt.Printf("user%02d", 11)
}
