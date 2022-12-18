package main

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/saintfish/chardet"
	"os"
)

func main() {
	file, err := os.ReadFile("360.log")
	if err != nil {
		fmt.Println(err)
	}
	detector := chardet.NewTextDetector()
	charset, err := detector.DetectBest(file)
	fmt.Println(charset.Charset)

	fileContent := string(file)
	decoder := mahonia.NewDecoder(charset.Charset)
	fileContent = decoder.ConvertString(fileContent)
	fmt.Println(fileContent)
}
