package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	src := `package main; func main(){fmt.Println("Hello, World!")}`

	// 设置编译参数
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 使用go/ast包遍历语法树
	ast.Inspect(f, func(n ast.Node) bool {
		fmt.Println(n)
		return true
	})
}
