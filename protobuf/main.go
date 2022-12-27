package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
)

func main() {
	// 初始化proto中的消息
	studyInfo := &StudyInfo{}

	//常规赋值
	studyInfo.Id = 1
	studyInfo.Name = "学习ProtoBuf"
	studyInfo.Duration = 180

	//在go中声明实例化map赋值给ProtoBuf消息中定义的map
	score := make(map[string]int32)
	score["实战"] = 100
	studyInfo.Score = score

	//用字符串的方式：打印ProtoBuf消息
	fmt.Printf("字符串输出结果：%v\n", studyInfo.String())

	//转成二进制文件
	marshal, err := proto.Marshal(studyInfo)
	if err != nil {
		return
	}
	fmt.Printf("Marshal转成二进制结果：%v\n", marshal)

	//将二进制文件转成结构体
	newStudyInfo := StudyInfo{}
	err = proto.Unmarshal(marshal, &newStudyInfo)
	if err != nil {
		return
	}
	fmt.Printf("二进制转成结构体的结果：%v\n", &newStudyInfo)

}
