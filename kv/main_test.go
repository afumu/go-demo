package main

import (
	"bufio"
	"fmt"
	dio "github.com/aquasecurity/go-dep-parser/pkg/io"
	"io"
	"os"
	"strings"
	"testing"
)

func TestGetConfigFromFile(t *testing.T) {
	open, _ := os.Open("os-release")
	releaseFileMap := GetConfigFromFile(open)

	osName := releaseFileMap["NAME"]
	osVersion := releaseFileMap["VERSION_ID"]
	fmt.Println(osName)
	fmt.Println(osVersion)
	if osName == "Kylin" && strings.HasPrefix(osVersion, "4") {
		fmt.Println("1111")
	}
	fmt.Printf("%+v", releaseFileMap)
}

func GetConfigFromFile(f dio.ReadSeekerAt) map[string]string {
	config := make(map[string]string)

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		config[key] = strings.Trim(value, "\"")
	}
	return config
}
