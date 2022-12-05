package main

import (
	"fmt"
	"regexp"
	"strings"
)

type User struct {
	Name string
}

func setUser(list []*User) {
	for _, v := range list {
		v.Name = "lisi"
	}
}

func main() {
	// 如果是upgrade文件名，不需要更新，防止upgrade二进制包错误导致永远无法更新
	var name = "-name=123    -age=123"
	name = DeleteExtraSpace(name)
	split := strings.Split(name, " ")
	for _, v := range split {
		fmt.Println(v)
	}
}

func DeleteExtraSpace(s string) string {
	regStr := "\\s{2,}"
	reg, _ := regexp.Compile(regStr)
	tmpStr := make([]byte, len(s))
	copy(tmpStr, s)
	spcIndex := reg.FindStringIndex(string(tmpStr))
	for len(spcIndex) > 0 {
		tmpStr = append(tmpStr[:spcIndex[0]+1], tmpStr[spcIndex[1]:]...)
		spcIndex = reg.FindStringIndex(string(tmpStr))
	}
	return string(tmpStr)
}
