package main

import (
	"net"
	"log"
	"io"
)

func main() {
	laddr := "127.0.0.1:80"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", laddr)
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
		return
	}

	const msgLen = 512
	for {
		c, err := ln.Accept()
		if err != nil {
			break
		}
		// Server connection.
		go func(c net.Conn) {
			defer func() {
				c.Close()
			}()
			//c.SetDeadline(time.Now().Add(time.Second * 3)) // Not intended to fire.



			for {

				var buf = make([]byte, 512)
				n, err := c.Read(buf)

				if err == io.EOF {
					log.Print("client EOL")
					break
				}else if err != nil {
					log.Fatal(err)
				}

				log.Print(string(buf[:n]))

				c.Write([]byte(`HTTP/1.1 200 OK
Date: Sat, 06 May 2017 07:27:23 GMT
Content-Type: text/html;charset=utf-8
Content-Length: 2

ok`))
			}

		}(c)
	}

}
