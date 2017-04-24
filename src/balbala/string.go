package main

import (
	"fmt"
	"os"
	"container/list"
	"reflect"
)

func main() {


	//stringcom()
	//calc()
	//array()
	//mymap()
	//ex()
	//mytype()
	//slice()
	mylist()
}
/*
链表及字节数组
反射获取类型
 */
func mylist(){
	l := list.New()

	bs := []byte{'a', 'b', 'c'}


	fmt.Print(reflect.TypeOf(string(bs)))

	fmt.Print(l.Back().Value.(string)+"\n")

	fmt.Print(string(bs)+"\n")

}

/*
动态数组（slice）
 */
func slice(){
	lines := []string{"1","2"}
	lines = append(lines,"3")
	fmt.Print(lines);
}

/*
字符串连接
 */
func stringcom() {
	s := "hello"
	s = "c" + s[1:] // 字符串虽不能更改，但可进行切片操作

	fmt.Printf("string: %s\n", s)
}

/*
数值计算
 */
func calc(){
	a, b, c := 1, 2, 3

	n := a + b + c;

	fmt.Printf("calc value: %s\n", n);
}

/*
数组操作
 */
func array(){
	a := [...]int{0,1,2,3,4,5,6,7,8}


	_ = a[3:4]

	fmt.Print(len(a));
}

func mymap(){
	//var numbers map[string]int
	// 另一种map的声明方式
	numbers := make(map[string]int)
	numbers["one"] = 1  //赋值
	numbers["ten"] = 10 //赋值
	numbers["three"] = 3

	fmt.Println("第三个数字是: ", numbers["three"]) // 读取数据

	for k, _ := range numbers{
		fmt.Print(k)
	}

}


func ex(){
	defer func() {
		if err := recover(); err != nil {
			_ = true
		}
	}()
	ext()
}

func ext() {
	var user = os.Getenv("USER")
	if user == "" {
		panic("no value for $USER")
	}
}

func mytype(){
	type person struct {
		name string
		age int
	}

	p := person{"abc", 18}

	fmt.Print(p)

	p2 := person{"def", 19}

	fmt.Print(p2)

}
