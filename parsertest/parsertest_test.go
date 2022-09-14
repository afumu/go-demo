package main

import (
	"fmt"
	"github.com/aquasecurity/go-dep-parser/pkg/golang/binary"
	"github.com/aquasecurity/go-dep-parser/pkg/java/jar"
	rpnVersion "github.com/knqyf263/go-rpm-version"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
	version "github.com/masahiro331/go-mvn-version"
	"os"
	"testing"
)

// 测试扫描jar包
func Test_jar(t *testing.T) {
	//var filePath = "D:\\doc\\log4j-core\\2.0-beta6\\log4j-core-2.0-beta6.jar"
	var filePath = "D:\\workplace\\junan-template\\target\\junan-template-1.0-SNAPSHOT.jar"
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

// 测试扫描go的二进制包
func TestGO(t *testing.T) {
	p := binary.NewParser()
	file, err := os.Open("./parsertest.exe")
	if err != nil {
		fmt.Println(err)
	}
	libs, deps, err := p.Parse(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(libs)
	fmt.Println(deps)
}

// 测试扫描rpm包
func TestRpm(t *testing.T) {
	db, err := rpmdb.Open("./Packages")
	if err != nil {
		fmt.Println(err)
	}

	pkgList, err := db.ListPackages()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pkgList)

}

// 测试java版本对比
func TestJavaVersion(t *testing.T) {
	v1, err := version.NewVersion("2.0-alpha")
	if err != nil {
		fmt.Println(err)
	}
	v2, err := version.NewVersion("2.0-beta9")
	if err != nil {
		fmt.Println(err)
	}

	comparer, _ := version.NewComparer("[2.0,2.8.2)")
	fmt.Println(comparer)

	fmt.Println(v1.GreaterThan(v2))

}

// 测试rpm包版本比较
func TestRpmVersion(t *testing.T) {
	v1 := rpnVersion.NewVersion("2:6.0-1")
	v2 := rpnVersion.NewVersion("2:6.0-2.el6")
	fmt.Println(v1.LessThan(v2))
}
