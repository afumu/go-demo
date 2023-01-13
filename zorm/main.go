package main

import (
	"context"
	"fmt"
	"gitee.com/chunanyong/zorm"
	"time"
)

type SysRole struct {
	Id          string `column:"ID"`
	RoleName    string `column:"ROLENAME"`
	RoleCode    string
	Description string
	CreateBy    string
	CreateTime  time.Time
	UpdateBy    string
	UpdateTime  time.Time
}

func (entity *SysRole) GetTableName() string {
	return "SYS_ROLE"
}

func (entity *SysRole) GetPKColumnName() string {
	//如果没有主键
	//return ""
	return "ID"
}

func newRoleStruct() SysRole {
	demo := SysRole{
		//如果Id=="",保存时zorm会调用zorm.FuncGenerateStringID(),默认时间戳+随机数,也可以自己定义实现方式,例如 zorm.FuncGenerateStringID=funcmyId
		Id:         zorm.FuncGenerateStringID(context.Background()),
		RoleName:   "defaultUserName",
		RoleCode:   "defaultPassword",
		CreateTime: time.Now(),
	}
	return demo
}

func main() {

	//声明一个对象的指针,用于承载返回的数据
	demo := &SysRole{}

	//构造查询用的finder
	finder := zorm.NewSelectFinder("SYS_ROLE") // select * from t_demo
	//finder = zorm.NewSelectFinder(demoStructTableName, "id,userName") // select id,userName from t_demo
	//finder = zorm.NewFinder().Append("SELECT * FROM " + demoStructTableName) // select * from t_demo
	//finder默认启用了sql注入检查,禁止语句中拼接 ' 单引号,可以设置 finder.InjectionCheck = false 解开限制

	//finder.Append 第一个参数是语句,后面的参数是对应的值,值的顺序要正确.语句统一使用?,zorm会处理数据库的差异
	finder.Append("WHERE role_code = ? ", "SysAdmin")

	//执行查询,has为true表示数据库有数据
	has, err := zorm.QueryRow(ctx, finder, demo)

	if err != nil { //标记测试失败
		fmt.Println(err)
		return
	}
	fmt.Println(has)
	//打印结果
	fmt.Println(demo)
}
