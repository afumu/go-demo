package main

import (
	"fmt"
	"github.com/aquasecurity/go-dep-parser/pkg/java/jar"
	"os"
)

func main() {
	jarDemo()
}

func jarDemo() {
	//var filePath = "D:\\doc\\log4j-core\\2.0-beta6\\log4j-core-2.0-beta6.jar"
	//var filePath = "E:\\vuln\\language\\java\\log4j-core\\2.0-beta6\\log4j-core-2.0-beta6.jar"
	var filePath = "E:\\vuln\\language\\java\\test\\ide-eval-resetter-2.2.3.jar"
	file, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err.Error())
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err.Error())
	}

	p := jar.NewParser(jar.WithSize(fileInfo.Size()), jar.WithFilePath(filePath), jar.WithOffline(true))

	libs, deps, err := p.Parse(file)

	fmt.Printf("%+v", libs)
	fmt.Printf("%+v", deps)
}
