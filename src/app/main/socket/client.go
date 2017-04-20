package main

import (
	"fmt"
	"net"
	"bufio"
)

func main() {
	proto := "HTTP/1.1"

	path := "/"

	method := "GET"

	addr := "127.0.0.1:80"

	eol := "\r\n"

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", addr)

	conn, _ := net.DialTCP("tcp", nil, tcpAddr)



	buf := bufio.NewWriter(conn);

	buf.WriteString(method + path + proto + eol)

	buf.WriteString(eol)

	buf_r := make([]byte, 2014)

	conn.Read(buf_r)

	fmt.Print("over")
}
