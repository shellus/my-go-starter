package main

import (
	"net"
	"github.com/xtaci/kcp-go"
	"io"
	"github.com/shellus/pkg/logs"
	"context"
)

type Tunnel struct {
	conn net.Conn
	kcp  net.Conn
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

	logs.Info("handleConn %s", conn.RemoteAddr().String())
	defer conn.Close()
	defer logs.Info("Conn close %s", conn.RemoteAddr().String())

	tunn, err := NewTunnel(conn)
	if err != nil {
		panic(err)
	}
	defer tunn.Close()
	tunn.Work()
}

func NewTunnel(conn net.Conn) (*Tunnel, error){
	kcpConn, err := kcp.DialWithOptions("127.0.0.1:8080", nil, 10, 3)
	if err != nil {
		return nil, err
	}
	tun := &Tunnel{
		conn:conn,
		kcp:kcpConn,
	}
	return tun, nil
}
func (tun *Tunnel) Work() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tun.cp(ctx)
}

func (tun *Tunnel) Close(){
	tun.kcp.Close()
}

func (tun *Tunnel) cp(parentCtx context.Context) {
	defer tun.kcp.Close()
	defer tun.conn.Close()

	ctx, cancel := context.WithCancel(parentCtx)

	go func() {
		_, err := io.Copy(tun.kcp, tun.conn);
		if err != nil {
			logs.Debug(err)
		}
		cancel()
	}()
	go func() {
		_, err := io.Copy(tun.conn, tun.kcp);
		if err != nil {
			logs.Debug(err)
		}
		cancel()
	}()
	<-ctx.Done()
	if err := ctx.Err(); err != nil && context.Canceled != err {
		logs.Error(err)
	}
}