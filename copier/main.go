package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"time"
)

type SysUser struct {
	Id        int64     // 标识
	Username  string    // 用户账号
	Realname  string    // 真实姓名
	Password  string    // 密码
	Salt      string    // md5密码盐
	Sex       int64     // 性别【0:未知，1-男，2-女】
	Avatar    string    // 头像
	Status    int64     // 状态【0-正常,1-冻结】
	DeletedAt time.Time // 删除状态【0-正常，1-已删除】
	CreatedBy int64     // 创建人
	CreatedAt time.Time // 创建时间
	UpdatedBy int64     // 更新人
	UpdatedAt time.Time // 更新时间
}

func main() {

	var userBefores []*SysUser
	start := time.Now()
	for i := 0; i < 100000; i++ {
		username := fmt.Sprintf("%s-%d", "用户账号", i)
		realname := fmt.Sprintf("%s-%d", "真实姓名", i)
		user := SysUser{Id: int64(i), Username: username, Realname: realname, Password: "asdfffffffffdf"}
		user.Salt = "asdfffffdfasdf"
		user.Sex = 0
		user.Avatar = "dfasdlfkjasdlkfjaklsdjflkasdjfklasdjflkajsdklfjasdklfja"
		user.Status = 0
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		userBefores = append(userBefores, &user)
	}
	end := time.Now()

	var userEnds []*SysUser
	start2 := time.Now()
	for _, v := range userBefores {
		// 1.
		user := SysUser{}
		copier.Copy(&user, &v)

		// 2.
		//user := SysUser{Id: v.Id, Username: v.Username, Realname: v.Realname, Password: v.Password}
		//user.Salt = v.Salt
		//user.Sex = v.Sex
		//user.Avatar = v.Avatar
		//user.Status = v.Status
		//user.CreatedAt = v.CreatedAt
		//user.UpdatedAt = v.UpdatedAt
		//userBefores = append(userBefores, &user)

		// 3.
		//user := copy[SysUser, SysUser](v)
		//
		userEnds = append(userEnds, &user)
	}
	end2 := time.Now()

	result := (end.UnixMilli() - start.UnixMilli())
	result2 := (end2.UnixMilli() - start2.UnixMilli())
	sprintf := fmt.Sprintf("耗时:%d 毫秒", result)
	sprintf2 := fmt.Sprintf("耗时:%d 毫秒", result2)
	fmt.Println(sprintf)
	fmt.Println(sprintf2)
}

func copy[T any, R any](t *T) *R {
	r := new(R)
	marshal, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(marshal, &r)
	if err != nil {
		fmt.Println(err)
	}
	return r
}
