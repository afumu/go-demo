package dao

import "fmt"

type User struct {
	ID         int64
	Username   string `gorm:"column:username"`
	Password   string `gorm:"column:password"`
	CreateTime int64  `gorm:"column:createtime"`
}

func (u User) TableName() string {
	return "users"
}

func SaveUser(user *User) {
	err := GlobalDb.Create(user).Error
	if err != nil {
		fmt.Println(err)
	}
}

func GetUserById() User {
	var u User
	err := GlobalDb.Where("id = ?", 1).Find(&u).Error
	if err != nil {
		fmt.Println(err)
	}
	return u
}

func GetAll() []User {
	var u []User
	err := GlobalDb.Find(&u).Error
	if err != nil {
		fmt.Println(err)
	}
	return u
}

func UpdateById() User {
	err := GlobalDb.Model(&User{}).Where("id = ?", 1).Update("username", "lisi").Error
	if err != nil {
		fmt.Println(err)
	}
	user := GetUserById()
	return user
}

func DeleteById() {
	err := GlobalDb.Where("id = ?", 1).Delete(&User{}).Error
	if err != nil {
		fmt.Println(err)
	}
}
