package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tangpanqing/aorm"
	"github.com/tangpanqing/aorm/base"
	"github.com/tangpanqing/aorm/builder"
	"github.com/tangpanqing/aorm/driver"
	"github.com/tangpanqing/aorm/null"
	"log"
	"time"
)

// 连接数据库
var db *base.Db

func init() {
	var err error
	db, err = aorm.Open(driver.Mysql, "root:root-abcd-1234@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
}

// 定义人员结构体
type Person struct {
	Id         null.Int    `aorm:"primary;auto_increment" json:"id"`
	Name       null.String `aorm:"size:100;not null;comment:名字" json:"name"`
	Sex        null.Bool   `aorm:"index;comment:性别" json:"sex"`
	Age        null.Int    `aorm:"index;comment:年龄" json:"age"`
	Type       null.Int    `aorm:"index;comment:类型" json:"type"`
	CreateTime null.Time   `aorm:"comment:创建时间" json:"createTime"`
	Money      null.Float  `aorm:"comment:金额" json:"money"`
	Test       null.Float  `aorm:"type:double;comment:测试" json:"test"`
}

// 定义文章结构体
type Article struct {
	Id          null.Int    `aorm:"primary;auto_increment" json:"id"`
	Type        null.Int    `aorm:"index;comment:类型" json:"type"`
	PersonId    null.Int    `aorm:"comment:人员Id" json:"personId"`
	ArticleBody null.String `aorm:"type:text;comment:文章内容" json:"articleBody"`
}

func main() {
	//quickStart()
	//实例化结构体
	var person = Person{}
	var article = Article{}

	//保存实例
	aorm.Store(&person, &article)
}

func quickStart() {
	//实例化结构体
	var person = Person{}
	var article = Article{}

	//保存实例
	aorm.Store(&person, &article)

	//迁移数据结构
	aorm.Migrator(db).AutoMigrate(&person, &article)

	//插入一个人员
	personId, _ := aorm.Db(db).Insert(&Person{
		Name:       null.StringFrom("Alice"),
		Sex:        null.BoolFrom(true),
		Age:        null.IntFrom(18),
		Type:       null.IntFrom(0),
		CreateTime: null.TimeFrom(time.Now()),
		Money:      null.FloatFrom(1),
		Test:       null.FloatFrom(2),
	})

	//插入一个文章
	articleId, _ := aorm.Db(db).Insert(&Article{
		Type:        null.IntFrom(0),
		PersonId:    null.IntFrom(personId),
		ArticleBody: null.StringFrom("文章内容"),
	})
	fmt.Println(articleId)

	//获取一条记录
	var personItem Person
	personColumn := Person{}
	err := aorm.Db(db).Table(&person).WhereEq(&personColumn.Name, "zouchangfu").OrderBy(&person.Id, builder.Desc).GetOne(&personItem)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("personItem:%+v", personItem)

	//联合查询
	/*	var list2 []ArticleVO
		aorm.
			Db(db).
			Table(&article).
			LeftJoin(&person, []builder.JoinCondition{
				builder.GenJoinCondition(&person.Id, builder.RawEq, &article.PersonId),
			}).
			SelectAll(&article).SelectAs(&person.Name, &articleVO.PersonName).
			WhereEq(&article.Id, articleId).
			GetMany(&list2)

		//带别名的联合查询
		var list3 []ArticleVO
		aorm.
			Db(db).
			Table(&article, "o").
			LeftJoin(&person, []builder.JoinCondition{
				builder.GenJoinCondition(&person.Id, builder.RawEq, &article.PersonId, "o"),
			}, "p").
			Select("*", "o").SelectAs(&person.Name, &articleVO.PersonName, "p").
			WhereEq(&article.Id, articleId, "o").
			GetMany(&list3)*/
}
