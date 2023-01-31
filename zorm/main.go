package main

import (
	"fmt"
	"gitee.com/chunanyong/zorm"
)

func main() {
	//initDm("dm://BSDS:bsds-abcd-1234@192.168.3.31:5236")
	initMysql("root:root-abcd-1234@tcp(192.168.3.31:3306)/bsds?charset=utf8&parseTime=true")
	finder := zorm.NewSelectFinder("sys_role")
	finder.Append("WHERE create_by = ? ", "admin")
	listMap, err := zorm.QueryMap(ctx, finder, nil)
	for _, v := range listMap {
		fmt.Println(v)
	}
	// 标记测试失败
	if err != nil {
		fmt.Println(err)
		return
	}
	

}
