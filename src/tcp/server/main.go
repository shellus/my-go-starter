package main

import (
	"net"
	"log"
	"fmt"
	"errors"
	"io"
	"strconv"
)

var errClose = errors.New("connect close")

func main() {
	laddr := "127.0.0.1:8081"
	tcpAddr, err := net.ResolveTCPAddr("tcp", laddr)
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("server init")

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go func(c net.Conn) {
			defer func() {
				c.Close()
				onClose(c)
				if r := recover(); r != nil {
					fmt.Printf("panic: %+v", r)
				}
			}()
			err = onConnect(c)

			if err != nil && err == errClose{
				return
			}

			for {
				var buf = make([]byte, 1024)
				n, err := c.Read(buf)
				if err != nil {
					if err == io.EOF {
						// 连接断开？
						break
					}else {
						panic(err)
					}
				}

				err = onRecv(c, buf[:n])

				if err != nil{
					if err == errClose {
						break
					}else {
						panic(err)
					}
				}

			}

		}(c)
	}

}

var arr = make(map [net.Conn]int)

func onConnect(c net.Conn)(err error) {
	arr[c] = 0
	fmt.Println("onConnect")
	return nil
}
func onClose(c net.Conn) {
	delete(arr, c)
	fmt.Println("onClose")
	return
}
func onRecv(c net.Conn, data []byte)(err error) {
	arr[c] += 1
	fmt.Println("onRecv")
	//fmt.Println(string(data))

	body := "ok:" + strconv.Itoa(arr[c])
	c.Write([]byte(fmt.Sprintf(`HTTP/1.1 200 OK
Date: Sat, 06 May 2017 07:27:23 GMT
Connection: keep-alive
Content-Type: text/html;charset=utf-8
Content-Length: %d

%s`, len(body), body)))
	//return errClose
	return nil
}