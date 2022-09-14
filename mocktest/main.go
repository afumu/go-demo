package main

import (
	"test-demo/mocktest/equipment"
)

func main() {

	phone := equipment.NewIphone6s()
	xiaoMing := NewPerson("xiaoMing", phone)
	xiaoMing.DayLife()
}
