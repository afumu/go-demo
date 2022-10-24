package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
	"unicode/utf8"
)

/**
  以下为GO新手可能遇到的一些坑：仅仅记录自己可能存在疑惑的测试代码
  参考学习地址：https://github.com/0voice/Introduction-to-Golang/blob/main/Golang%20%E6%96%B0%E6%89%8B%E5%8F%AF%E8%83%BD%E4%BC%9A%E8%B8%A9%E7%9A%84%2050%20%E4%B8%AA%E5%9D%91.md
*/

// 我们在使用slice的时候是可以通过cap来判断分配空间的大小，但是map是不可以的
func TestMap(t *testing.T) {
	var arr []string
	arr = append(arr, "abc")
	arr = append(arr, "abc")
	arr = append(arr, "abc")
	arr = append(arr, "abc")
	arr = append(arr, "abc")
	fmt.Println(len(arr), cap(arr))

	m := make(map[string]int)
	m["a"] = 1
	fmt.Println(len(m))
	// 这里是编译不通过的
	//fmt.Println(cap(m))
}

// ---------------------------------------------------------------------------------------------------------------------
// 指针类型数组
// 在 C/C++ 中，数组（名）是指针。将数组作为参数传进函数时，相当于传递了数组内存地址的引用，在函数内部会改变该数组的值。
// 在 Go 中，数组是值。作为参数传进函数时，传递的是数组的原始值拷贝，此时在函数内部是无法更新该数组的：
func TestSlice(t *testing.T) {

	/*	x := [3]int{1, 2, 3}
		func(arr [3]int) {
			arr[0] = 7
			fmt.Println(arr) // [7 2 3]
		}(x)
		fmt.Println(x) // [1 2 3]    // 并不是你以为的 [7 2 3]*/

	/*	x := [3]int{1, 2, 3}
		func(arr *[3]int) {
			// 修改指针类型数组内部的数据
			(*arr)[0] = 7
			fmt.Println(arr) // &[7 2 3]
		}(&x)
		fmt.Println(x) // [7 2 3]*/

	// 如果是一个切片的话，就直接可以修改了
	x := []int{1, 2, 3}
	func(arr []int) {
		arr[0] = 7
		fmt.Println(x) // [7 2 3]
	}(x)
	fmt.Println(x) // [7 2 3]
}

// ---------------------------------------------------------------------------------------------------------------------
// Go 则会返回元素对应数据类型的零值，比如 nil、” 、false 和 0，取值操作总有值返回，故不能通过取出来的值来判断 key 是不是在 map 中。
// 检查 key 是否存在可以用 map 直接访问，检查返回的第二个参数即可：
// 检测map的key是否存在
func TestMapKeyCheck(t *testing.T) {
	// 这是一个错误的检测key是否存在
	/*	x := map[string]string{"one": "2", "two": "", "three": "3"}
		if v := x["two"]; v == "" {
			fmt.Println("key two is no entry") // 键 two 存不存在都会返回的空字符串
		}*/

	// 正确判断方式是通过第二个参数来判断
	x := map[string]string{"one": "2", "two": "", "three": "3"}
	if _, ok := x["two"]; !ok {
		fmt.Println("key two is no entry")
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// 尝试使用索引遍历字符串，来更新字符串中的个别字符，是不允许的。
// string 类型的值是只读的二进制 byte slice，
// 如果真要修改字符串中的字符，将 string 转为 []byte 修改后，再转为 string 即可：
func TestString(t *testing.T) {
	// string 类型不可以更改
	//x := "text"
	//x[0] = "T"        // error: cannot assign to x[0]
	//fmt.Println(x)

	/*	x := "text"
		xBytes := []byte(x)
		fmt.Println(xBytes)
		xBytes[0] = 'T' // 注意此时的 T 是 rune 类型
		x = string(xBytes)
		fmt.Println(x) // Text*/

	// 上面的做法其实也是有问题的，一个中文可能占据好几个字节
	x := "text"
	xRunes := []rune(x)
	xRunes[0] = '我'
	x = string(xRunes)
	fmt.Println(x) // 我ext
}

// ---------------------------------------------------------------------------------------------------------------------
// 对字符串用索引访问返回的不是字符，而是一个 byte 值。
func TestStringIndex(t *testing.T) {
	x := "ascii"
	fmt.Println(x[0])        // 97
	fmt.Printf("%T\n", x[0]) // uint8
}

// ---------------------------------------------------------------------------------------------------------------------
// string 的值不必是 UTF8 文本，可以包含任意的值。只有字符串是文字字面值时才是 UTF8 文本，字串可以通过转义来包含其他数据。
// 判断字符串是否是 UTF8 文本，可使用 "unicode/utf8" 包中的 ValidString() 函数：
func TestStringUtf8(t *testing.T) {
	str1 := "ABC"
	fmt.Println(utf8.ValidString(str1)) // true

	str2 := "A\xfeC"
	fmt.Println(utf8.ValidString(str2)) // false

	str3 := "A\\xfeC"
	fmt.Println(utf8.ValidString(str3)) // true    // 把转义字符转义成字面值
}

// ---------------------------------------------------------------------------------------------------------------------
// 字符串的长度
func TestStringLength(t *testing.T) {
	char := "♥"
	fmt.Println(len(char)) // 3

	fmt.Println(utf8.RuneCountInString(char)) // 1
	fmt.Println(len([]rune(char)))            // 1

	fmt.Println("------------------------------------")
	char2 := "é"
	fmt.Println(len(char2))                    // 2
	fmt.Println(utf8.RuneCountInString(char2)) // 1
	fmt.Println("cafe\u0301")                  // café    // 法文的 cafe，实际上是两个 rune 的组合
}

// ---------------------------------------------------------------------------------------------------------------------
// 范围查询遍历字符
func TestRangString(t *testing.T) {
	data := "A\xfe\x02\xff\x04"
	for _, v := range data {
		fmt.Printf("%#x ", v) // 0x41 0xfffd 0x2 0xfffd 0x4    // 错误
	}

	fmt.Println("---------------------------")
	for _, v := range []byte(data) {
		fmt.Printf("%#x ", v) // 0x41 0xfe 0x2 0xff 0x4    // 正确
	}
}

func TestWaitGo(t *testing.T) {
	var wg sync.WaitGroup
	done := make(chan struct{})

	workerCount := 2
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go doIt(i, done, &wg)
	}

	close(done)
	wg.Wait()
	fmt.Println("all done!")
}

func doIt(workerID int, done <-chan struct{}, wg *sync.WaitGroup) {
	fmt.Printf("[%v] is running\n", workerID)
	defer wg.Done()
	<-done
	fmt.Printf("[%v] is done\n", workerID)
}

// ---------------------------------------------------------------------------------------------------------------------
// 通知剩余的协程不需要工作了
func TestName(t *testing.T) {
	ch := make(chan int)
	done := make(chan struct{})

	for i := 0; i < 3; i++ {
		go func(idx int) {
			select {
			case ch <- (idx + 1) * 2:
				fmt.Println(idx, "Send result")
			case <-done:
				fmt.Println(idx, "Exiting")
			}
		}(i)
	}

	fmt.Println("Result: ", <-ch)
	close(done)
	time.Sleep(3 * time.Second)
}

// ---------------------------------------------------------------------------------------------------------------------
// 测试动态关闭channel
func TestChannelKill(t *testing.T) {
	inCh := make(chan int)
	outCh := make(chan int)

	go func() {
		var in <-chan int = inCh
		var out chan<- int
		var val int

		for {
			select {
			case out <- val:
				println("--------")
				out = nil
				in = inCh
			case val = <-in:
				println("++++++++++")
				out = outCh
				in = nil
			}
		}
	}()

	go func() {
		for r := range outCh {
			fmt.Println("Result: ", r)
		}
	}()

	time.Sleep(0)
	inCh <- 1
	inCh <- 2
	time.Sleep(3 * time.Second)
}

func TestHttp(t *testing.T) {
	resp, err := http.Get("https://api.ipify.org1/?format=json111")
	fmt.Println(resp)
	if err != nil {
		fmt.Println("err:", err)
	}
	if resp != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println("err:", err)
			}
		}(resp.Body)
	}
}
