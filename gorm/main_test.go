package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

/*
	type Product struct {
		gorm.Model
		Code  string
		Price uint
	}

	func Test_Gorm(t *testing.T) {
		db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		// 迁移 schema
		db.AutoMigrate(&Product{})

		// Create
		db.Create(&Product{Code: "D42", Price: 100})

		// Read
		var product Product
		db.First(&product, 1)                 // 根据整型主键查找
		db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

		// Update - 将 product 的 price 更新为 200
		db.Model(&product).Update("Price", 200)
		// Update - 更新多个字段
		db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
		db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

		// Delete - 删除 product
		db.Delete(&product, 1)
	}
*/

type Test1 struct {
	gorm.Model
	Code  string
	Price uint
}

var globalDb *gorm.DB

func init() {
	dsn := "root:root-abcd-1234@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	globalDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
}

func TestCreate(t *testing.T) {
	globalDb.AutoMigrate(&Test1{})
	// Create
	test1 := Test1{Code: "D43", Price: 100}
	globalDb.Create(&test1)
}

func TestFind(t *testing.T) {
	// Read
	var test1 Test1
	globalDb.First(&test1, 1)                 // 根据整型主键查找
	globalDb.First(&test1, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	fmt.Println(test1)
}

func TestUpdate(t *testing.T) {
	var test1 Test1
	globalDb.First(&test1, 1)
	globalDb.First(&test1, "code = ?", "D42")

	// Update - 将 product 的 price 更新为 200
	globalDb.Model(&test1).Update("Price", 200)

	// Update - 更新多个字段
	globalDb.Model(&test1).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	globalDb.Model(&test1).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

}

func TestFindFirst(t *testing.T) {
	// Read
	var test1 Test1
	globalDb.First(&test1, 1)                 // 根据整型主键查找
	globalDb.First(&test1, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	fmt.Println(test1)
}

func TestFindTake(t *testing.T) {
	// Read
	var test1 Test1
	globalDb.Take(&test1, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	fmt.Println(test1)
}

func TestFindLast(t *testing.T) {
	// Read
	var test1 Test1
	globalDb.Last(&test1, "code = ?", "D42")
	fmt.Println(test1)
}
