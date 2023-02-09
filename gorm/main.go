package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

var gormDb *gorm.DB

type Student struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	Age       uint8
	Email     string
	Birthday  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func init() {
	dsn := "root:root-abcd-1234@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	gormDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println(err)
	}
}

func main() {

}

func funcName() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	type User struct {
	}

	var users []*User
	sqlResult := db.Where("username = ? and age = ?", "zhangsan", 18).Find(&users)
	fmt.Println(sqlResult)

	/*	// 迁移 schema
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
		db.Delete(&product, 1)*/
}
