package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	args := os.Args
	fmt.Println(args)
	paramMap := make(map[string]interface{})
	paramMap["name"] = "张三"
	paramMap["age"] = 123
	params, _ := convertParams("http://192.168.3.31:80", paramMap)
	fmt.Println(params)
}

func convertParams(reqUrl string, paramMap map[string]interface{}) (string, error) {
	params := url.Values{}
	for k, v := range paramMap {
		params.Set(k, fmt.Sprintf("%v", v))
	}
	u, err := url.ParseRequestURI(reqUrl)
	if err != nil {
		return "", err
	}
	u.RawQuery = params.Encode()
	return u.String(), err
}
