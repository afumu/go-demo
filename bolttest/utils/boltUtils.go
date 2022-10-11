package utils

import (
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
)

var bucket *bolt.Bucket

func init() {
	var result string
	if runtime.GOOS == "windows" {
		result = Command("echo", "%USERPROFILE%")
	} else {
		result = Command("whoami")
	}

	path := fmt.Sprintf("%s%s.situation-awareness-data%sdb%sseek.db",
		Trim(result), string(filepath.Separator), string(filepath.Separator), string(filepath.Separator))
	var err error
	globalDb, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Println("init db error：", err)
	}

	err = globalDb.Update(func(tx *bolt.Tx) error {
		bucket, err = tx.CreateBucketIfNotExists([]byte("default"))
		return err
	})

	if err != nil {
		log.Println("init db error：", err)
	}
}

func Trim(src string) (dist string) {
	if len(src) == 0 {
		return
	}
	var distR []rune
	r := []rune(src)
	fmt.Println(r)
	for i := 0; i < len(r); i++ {
		if r[i] == 10 || r[i] == 13 || r[i] == 32 {
			continue
		}
		distR = append(distR, r[i])
	}
	return string(distR)
}

func Command(arg ...string) (result string) {
	name := "/bin/bash"
	c := "-c"

	// 根据系统设定不同的命令name
	if runtime.GOOS == "windows" {
		name = "cmd"
		c = "/C"
	}

	arg = append([]string{c}, arg...)
	cmd := exec.Command(name, arg...)

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		log.Printf("Error:The command is err,%s", err)
		return
	}

	//读取所有输出
	bytesResult, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Printf("ReadAll Stdout:%s", err.Error())
		return
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("wait:%s", err.Error())
		return
	}

	return string(bytesResult)
}

func Get(key string) {
	bucket.Get([]byte(key))
}

func Put(key string, v string) error {
	return bucket.Put([]byte(key), []byte(v))
}
