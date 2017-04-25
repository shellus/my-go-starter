package main

import (
	"strings"
	"unsafe"
	"fmt"
)
// 定义一个和 strings 包中的 Reader 相同的本地结构体
type Reader struct {
	s        string
	i        int64
	prevRune int
}

func main() {
	// 创建一个 strings 包中的 Reader 对象
	sr := strings.NewReader("abcdef")
	// 此时 sr 中的成员是无法修改的
	fmt.Println(sr)
	// 我们可以通过 unsafe 来进行修改
	// 先将其转换为通用指针
	p := unsafe.Pointer(sr)
	// 再转换为本地 Reader 结构体
	pR := (*Reader)(p)
	// 这样就可以自由修改 sr 中的私有成员了
	(*pR).i = 3 // 修改索引
	// 看看修改结果
	fmt.Println(sr)
	// 看看读出的是什么
	b, err := sr.ReadByte()
	fmt.Printf("%c, %v\n", b, err)
}