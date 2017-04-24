package main

import (
	"net"
	"fmt"
	"bufio"
	"strconv"
)

const PORT = 3540

func main() {
	server, err := net.Listen("tcp", ":" + strconv.Itoa(PORT))
	if server == nil {
		panic(err)
	}
	conns := clientConns(server)
	for {
		go handleConn(<-conns)
	}
}

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Errorf("couldn't accept: %s", err.Error())
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn) {
	b := bufio.NewReader(client)
	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			// EOF, or worse
			break
		}
		client.Write(line)
	}
}