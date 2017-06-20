package main

import (
	"github.com/tarm/goserial"
	"fmt"
)

func main() {
	c := &serial.Config{Name: "COM5", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}

	n, err := s.Write([]byte("u"))
	if err != nil {
		panic(err)
	}
	fmt.Println(n)

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q", buf[:n])
}
