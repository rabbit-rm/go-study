package main

import (
	"fmt"
)

/*func main() {
	type foo struct {
		bar string
	}
	s1 := []foo{
		{"A"},
		{"B"},
		{"C"},
	}
	var s2 = make([]*foo, 3)

	for i, v := range s1 {
		s2[i] = &v
	}
	fmt.Println(s1[0], s1[1], s1[2])
	fmt.Println(s2[0], s2[1], s2[2])
}*/

func main() {
	// go 1.22 之前，for-range 使用:=声明的变量在每次循环中都会被重用
	// go 1.22，语义改变，for-range 使用:=声明的变量在每次循环中都会被重新分配地址
	var arr = [3]int{}
	for i, value := range arr {
		fmt.Printf("i:%v,value:%v,&value:%p\n", i, value, &value)
	}
	// Output:
	// i:0,value:0,&value:0xc00000a0b8
	// i:1,value:0,&value:0xc00000a0d8
	// i:2,value:0,&value:0xc00000a100
}
