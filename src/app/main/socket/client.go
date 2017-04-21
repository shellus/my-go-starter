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

	headers := []string{
		"Host: 127.0.0.1",
		"Connection: keep-alive",
		"Cache-Control: max-age=0",
		"Upgrade-Insecure-Requests: 1",
		"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		//"Accept-Encoding: gzip, deflate, sdch, br",
		"Accept-Language: zh-CN,zh;q=0.8",
	}

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", addr)

	conn, _ := net.DialTCP("tcp", nil, tcpAddr)

	buf := bufio.NewWriter(conn);

	buf.WriteString(method + " " + path + " " + proto + eol)

	for i := 0; i < len(headers); i++ {
		buf.WriteString(headers[i] + eol)
	}
	buf.WriteString(eol)

	buf.Flush()

	buf_r := make([]byte, 2014)

	recv_len, _ := conn.Read(buf_r)

	fmt.Printf("over %d : %s", recv_len, buf_r)
}
