package main

import (
	"net"
	"github.com/xtaci/kcp-go"
	"io"
	"log"
	"context"
	"os"
	"fmt"
)

type Tunnel struct {
	clientConn net.Conn
	serverConn net.Conn
}

func init() {

}
func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {

	log.Printf("handleConn %s", conn.RemoteAddr().String())
	defer log.Printf("Conn close %s", conn.RemoteAddr().String())
	defer conn.Close()
	kcpconn, err := kcp.DialWithOptions("127.0.0.1:8080", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	defer kcpconn.Close()

	tunn := &Tunnel{
		clientConn:conn,
		serverConn:kcpconn,
	}
	tunn.Start()
}
func (tunn *Tunnel) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tunn.copy(ctx)
}

func (tunn *Tunnel) copy(ctx context.Context) {
	defer tunn.serverConn.Close()
	defer tunn.clientConn.Close()

	cctx, cancel := context.WithCancel(ctx)

	go func() {
		_, err := io.Copy(tunn.serverConn, tunn.clientConn);
		if err != nil {
			tunn.logger().Printf("%s", err)
		}
		cancel()
	}()
	go func() {
		_, err := io.Copy(tunn.clientConn, tunn.serverConn);
		if err != nil {
			tunn.logger().Printf("%s", err)
		}
		cancel()
	}()

	cctx.Done()
}

func (tunn *Tunnel) logger() *log.Logger {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	logger.SetPrefix(fmt.Sprintf("s=%s,c=%s | ", tunn.serverConn.RemoteAddr().String(), tunn.clientConn.RemoteAddr().String()))
	return logger
}