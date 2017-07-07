package main

import (
	"bytes"
	"github.com/shellus/pkg/sutil"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"github.com/shellus/my-go-starter/src/scanner/rand_reader"
)

type MyPackage struct {
	L1  uint32
	L2  uint32
	L3  uint32
	Buf []byte
}

func main(){
	buf := rand_reader.New([]byte{})
	buf.Write(generator())
	buf.Write(generator())
	buf.Write(generator())
	buf.Write(generator())

	fmt.Println()


	dec := gob.NewDecoder(buf)
	pack := new(MyPackage)

	for{
		err := dec.Decode(pack)
		if err != nil {
			panic(err)
		}
		fmt.Println(len(pack.Buf))
	}


}

func generator()[]byte{
	var randBytesLen int

	randBytesLen = sutil.RandInt(1024*100, 1024* 1024)

	randBytes := make([]byte, randBytesLen)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)
	}

	pack := &MyPackage{}
	pack.L1 =0
	pack.L2 =0
	pack.L3 =uint32(randBytesLen)
	fmt.Println(pack.L3)
	pack.Buf = randBytes

	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)

	err = enc.Encode(pack)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("randPackage len: %d\n", buf.Len())
	return buf.Bytes()
}