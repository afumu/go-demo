package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

// 解析 window 下的命令systeminfo 的信息
// 根据systeminfo的信息进行window下的漏洞扫描工作
func Test_ParseSystemInfo(t *testing.T) {
	info := getSystemInfo()

	// os version
	systeminfo_matches := GetValueStringByRegex(info, ".*?((\\d+\\.?){3}) ((Service Pack (\\d)|N\\/\\w|.+) )?[ -\\xa5]+ (\\d+).*")
	fmt.Printf("%+v\n", systeminfo_matches)

	// OS Name
	win_matches := GetValueStringByRegex(info, ".*?Microsoft[\\(R\\)]{0,3} Windows[\\(R\\)?]{0,3} ?(Serverr? )?(\\d+\\.?\\d?( R2)?|XP|VistaT).*")
	fmt.Printf("%+v\n", win_matches)

	//System Type
	archs := GetValueStringByRegex(info, ".*?([\\w\\d]+?)-based PC.*")
	fmt.Printf("%+v\n", archs)

	// 补丁
	compile := regexp.MustCompile(".*KB\\d+.*")
	all := compile.FindAll([]byte(info), -1)
	fmt.Println(all)
	for _, v := range all {
		kbs := GetValueStringByRegex(string(v), ".*KB(\\d+).*")
		fmt.Printf("%+v\n", kbs)
	}

}

func GetValueStringByRegex(str, rule string) []string {
	reg, err := regexp.Compile(rule)
	if reg == nil || err != nil {
		return nil
	}
	//提取关键信息
	result := reg.FindStringSubmatch(str)
	if len(result) < 2 {
		return nil
	}
	return result
}

// 获取window的系统信息
func getSystemInfo() string {
	cmd := exec.Command("systeminfo")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return ""
	}

	in := bufio.NewScanner(stdout)
	builder := strings.Builder{}
	for in.Scan() {
		cmdRe := ConvertByte2String(in.Bytes(), "GB18030")
		builder.Write([]byte(cmdRe + "\n"))
	}
	err = cmd.Wait()
	if err != nil {
		log.Println(err)
		return ""
	}
	systemInfo := builder.String()

	//fmt.Println(systemInfo)
	return systemInfo
}

// 解码
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
