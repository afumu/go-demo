package main

import (
	"fmt"
	"testing"

	. "github.com/dave/jennifer/jen"
)

var UserColumn = struct {
	ID        string
	Username  string
	Password  string
	Address   string
	Age       string
	Phone     string
	Score     string
	Dept      string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	Username:  "username",
	Password:  "password",
	Address:   "address",
	Age:       "age",
	Phone:     "phone",
	Score:     "score",
	Dept:      "dept",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}
var UserField = struct {
	ID string
}{}

func TestName(t *testing.T) {
	/*	f := NewFile("main")
		f.Func().Id("main").Params().Block(
			Qual("fmt", "Println").Call(Lit("Hello, world")),
		)*/

	f := Var().Id("UserField").Op("=").Id("struct").Id("{").
		Id("\n").Id("ID").String().Id("\n").Id("}").Id("{").Id("}")
	fmt.Printf("%#v\n", f)
}
