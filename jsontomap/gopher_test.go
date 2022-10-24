package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestMap(t *testing.T) {

	file, _ := os.ReadFile("1.json")

	var cmap = make(map[string]interface{})

	err := json.Unmarshal(file, &cmap)
	fmt.Println(err)
	fmt.Println(cmap)
	openPort := cmap["data"].(map[string]interface{})["openPort"]
	fileName := cmap["data"].(map[string]interface{})["fileName"]
	sensitiveService := cmap["data"].(map[string]interface{})["sensitiveService"]
	requiredApplication := cmap["data"].(map[string]interface{})["requiredApplication"]
	openPortStr := fmt.Sprintf("%s", openPort)
	fmt.Println("openPortStr:", openPortStr)
	fmt.Println(fileName)
	fmt.Println(sensitiveService)

	sprintf := fmt.Sprintf("%s", requiredApplication)
	fmt.Println(sprintf)

}

func print(string2 string) {
	fmt.Printf(string2, "bb")
}
