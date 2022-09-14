package equipment

import "fmt"

type Iphone6s struct {
}

func NewIphone6s() *Iphone6s {
	return &Iphone6s{}
}

func (p *Iphone6s) WeiXin() bool {
	fmt.Println("Iphone6s chat wei xin!")
	return true
}

func (p *Iphone6s) WangZhe() bool {
	fmt.Println("Iphone6s play wang zhe!")
	return true
}

func (p *Iphone6s) ZhiHu() bool {
	fmt.Println("Iphone6s read zhi hu!")
	return true
}
