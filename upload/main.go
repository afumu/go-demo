package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// 本地文件路径
	filePath := "C:\\Program Files\\agent\\bin\\log\\agent-2023-05-13T06-58-49.663.log.gz"

	// 创建一个buffer来存储文件内容
	fileBuffer := &bytes.Buffer{}
	fileWriter := multipart.NewWriter(fileBuffer)

	// 创建一个文件字段
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	fileWriter.WriteField("file", "agent-2023-05-13T06-58-49.663.log.gz") // 文件字段的名称和文件名

	// 将文件内容拷贝到buffer中
	filePart, err := fileWriter.CreateFormFile("file", "agent-2023-05-13T06-58-49.663.log.gz")
	if err != nil {
		fmt.Println("创建表单文件字段失败:", err)
		return
	}
	_, err = io.Copy(filePart, file)
	if err != nil {
		fmt.Println("拷贝文件内容失败:", err)
		return
	}

	// 完成写入并生成表单数据的Content-Type头部
	fileWriter.Close()

	// 创建一个HTTP请求
	request, err := http.NewRequest("POST", "http://192.168.3.104:81/sys/sysFile/upload", fileBuffer)
	if err != nil {
		fmt.Println("创建HTTP请求失败:", err)
		return
	}
	request.Header.Set("Content-Type", fileWriter.FormDataContentType())
	request.Header.Set("Cookie", "sl_langtype=zh; lang=zh-CN; X-Access-Token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2ODYwMzI2MDksInVzZXJuYW1lIjoiYWRtaW4ifQ.ysYAFI1sEvtkoI0sh8_SF7055b8D6tBu6GryO-Ykd3k")

	// 发送请求
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("发送HTTP请求失败:", err)
		return
	}
	defer response.Body.Close()

	// 处理响应
	if response.StatusCode == http.StatusOK {
		fmt.Println("文件上传成功")
	} else {
		fmt.Println("文件上传失败:", response.Status)
	}
}
