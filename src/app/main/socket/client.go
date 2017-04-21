package main

import (
	"fmt"
	"net"
	"bufio"
	"runtime"
)
type Request struct {
	addr string
	method string
	path string
	proto string
	headers []string
}
type Response struct{
	data []byte
	data_len int
}

func main() {
	request := Request{}

	request.proto = "HTTP/1.1"

	request.path = "/"

	request.method = "GET"

	request.addr = "127.0.0.1:80"

	request.headers = []string{
		"Host: 127.0.0.1",
		"Connection: keep-alive",
		"Cache-Control: max-age=0",
		"Upgrade-Insecure-Requests: 1",
		"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		//"Accept-Encoding: gzip, deflate, sdch, br",
		"Accept-Language: zh-CN,zh;q=0.8",
	}
	c := make(chan Response);
	for i := 0; i < 10; i++ {
		go http_request(request, c)
	}

	n:=0
	for response := range c {
		n++
		fmt.Printf("response.data_len %d %d \n", response.data_len, n)
	}


}

func http_request(request Request, c chan Response){
	eol := "\r\n"
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", request.addr)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	buf := bufio.NewWriter(conn);

	for err == nil {
		buf.WriteString(request.method + " " + request.path + " " + request.proto + eol)

		for i := 0; i < len(request.headers); i++ {
			buf.WriteString(request.headers[i] + eol)
		}
		buf.WriteString(eol)

		buf.Flush()

		buf_r := make([]byte, 2014)

		recv_len, _ := conn.Read(buf_r)

		c <- Response{data:buf_r, data_len:recv_len}

		runtime.Gosched()
	}
}
