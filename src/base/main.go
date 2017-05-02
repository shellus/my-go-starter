package main

import (
	"fmt"
	"os"
	"container/list"
	"reflect"
	"strconv"
	"queue"
	"time"
)

func main() {

	//stringcom()
	//calc()
	//array()
	//mymap()
	//ex()
	//mytype()
	//slice()
	//mylist()
	//convert()
	queue_test()
}
func queue_test(){

	// 并发2
	q := queue.NewQueue(2)

	// 订阅
	q.Sub(func(j queue.Job) {
		time.Sleep(time.Millisecond * 100)
		fmt.Println(j.Value.(string))
	})

	// 发布
	for i:=0; i<10 ;i++ {
		q.Push(queue.Job{Value:fmt.Sprintf("lalala: %d", i)})
	}

	// 开始工作
	q.Work()

}
/**
字符串，数值转换
 */
func convert() {

	myi, _ := strconv.Atoi("123")
	fmt.Printf("convert to int: %d \n", myi)

	fmt.Printf("convert to int string: %s \n", strconv.Itoa(123))

	bytes := []byte{'a', 'b', 'c'}
	fmt.Printf("convert to string: %s \n", string(bytes))

	mystring := "abc"
	fmt.Printf("convert to bytes: %s \n", []byte(mystring))

	b, _ := strconv.ParseBool("true")
	f, _ := strconv.ParseFloat("3.1415", 64)
	i, _ := strconv.ParseInt("-42", 10, 64)
	u, _ := strconv.ParseUint("42", 10, 64)

	fmt.Print(b, f, i, u)
}
/*
链表及字节数组
反射获取类型
 */
func mylist() {

	l := list.New()

	bs := []byte{'a', 'b', 'c'}

	fmt.Print(reflect.TypeOf(string(bs)))

	fmt.Print(l.Back().Value.(string) + "\n")

	fmt.Print(string(bs) + "\n")

}

/*
动态数组（slice）
 */
func slice() {
	lines := []string{"1", "2"}
	lines = append(lines, "3")
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
func calc() {
	a, b, c := 1, 2, 3

	n := a + b + c;

	fmt.Printf("calc value: %s\n", n);
}

/*
数组操作
 */
func array() {
	a := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8}

	_ = a[3:4]

	fmt.Print(len(a));
}

func mymap() {
	//var numbers map[string]int
	// 另一种map的声明方式
	numbers := make(map[string]int)
	numbers["one"] = 1  //赋值
	numbers["ten"] = 10 //赋值
	numbers["three"] = 3

	fmt.Println("第三个数字是: ", numbers["three"]) // 读取数据

	for k, _ := range numbers {
		fmt.Print(k)
	}

}

func ex() {
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

func mytype() {
	type person struct {
		name string
		age  int
	}

	p := person{"abc", 18}

	fmt.Print(p)

	p2 := person{"def", 19}

	fmt.Print(p2)

}
