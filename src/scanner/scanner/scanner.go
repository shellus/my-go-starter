package main

import (
	"encoding/binary"
	"bytes"
	"crypto/rand"
	"math/big"
	"fmt"
	"bufio"
	"github.com/shellus/my-go-starter/src/scanner/rand_reader"
)
var headerLen int = 4

func main(){
	// 随机生成一个包，并创建一个模拟conn的reader
	conn := rand_reader.New(randPackage())

	// 新建扫描器
	scanner := bufio.NewScanner(conn)

	// 设定分页规则
	split := func(data []byte, atEOF bool) (adv int, token []byte, err error) {
		//fmt.Println(len(data))
		// header还没接受完（不够12字节）
		if len(data) < headerLen*3 {
			return
		}
		var l1, l2, l3 uint32
		buf := bytes.NewReader(data)
		binary.Read(buf, binary.LittleEndian, &l1)
		binary.Read(buf, binary.LittleEndian, &l2)
		binary.Read(buf, binary.LittleEndian, &l3)

		// body不够长
		if len(data) < headerLen * 3 + int(l3) {
			return
		}

		// 够了，告诉解析器，这次前进多少n，和这次拿到的数据bytes
		return headerLen*3 + int(l3), data[:headerLen*3 + int(l3)], nil
	}
	scanner.Split(split)

	// 默认一次扫描最大64K，这里设置为1M
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		// 每次scanner.Scan()return true，代表扫描到一个包。
		if decodePackage(scanner.Bytes()) {
			fmt.Println("decodePackage success")
		}else {
			fmt.Println("decodePackage fail")
		}
	}

	// 扫描完了之后（scanner.Scan()返回false），检查扫描是否有错误
	if scanner.Err() != nil {
		panic(scanner.Err().Error())
	}
}

func decodePackage(b []byte) bool {
	//fmt.Printf("decodePackage len: %d\n", len(b))
	buf := bytes.NewReader(b)
	var l1,l2,l3 uint32
	binary.Read(buf, binary.LittleEndian, &l1)
	binary.Read(buf, binary.LittleEndian, &l2)
	binary.Read(buf, binary.LittleEndian, &l3)

	if len(b) == headerLen * 3 + int(l3) {
		return true
	}else {
		return false
	}
}
func randPackage() []byte{
	var randBytesLen int

	{
		i, err := rand.Int(rand.Reader, new(big.Int).SetInt64(int64(1024 * 100)))
		if err != nil {
			panic(err)
		}
		randBytesLen = int(i.Int64())
	}

	randBytes := make([]byte, randBytesLen)

	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = binary.Write(buf, binary.LittleEndian, uint32(0))
	err = binary.Write(buf, binary.LittleEndian, uint32(0))
	err = binary.Write(buf, binary.LittleEndian, uint32(randBytesLen))
	err = binary.Write(buf, binary.LittleEndian, randBytes)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("randPackage len: %d\n", buf.Len())
	return buf.Bytes()
}