package test1

import (
	"fmt"
	"reflect"
	"testing"
)

type Student struct {
	Name  string
	Age   int
	Score float64
}

// 直接使用 new(T) 创建指针对象的性能测试
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		student := new(Student)
		_ = student
	}
}

// 使用 reflect.New 创建指针对象的性能测试
func BenchmarkReflectNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		studentPtr := reflect.New(reflect.TypeOf((*Student)(nil)).Elem()).Interface()
		_ = studentPtr
	}
}

func main() {
	fmt.Println("Benchmark New:")
	testing.Benchmark(BenchmarkNew)
	fmt.Println("Benchmark Reflect New:")
	testing.Benchmark(BenchmarkReflectNew)
}
