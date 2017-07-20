package main

import (
	"github.com/armon/go-socks5"
	"github.com/shellus/pkg/logs"
	"flag"
)

func main(){
	var listenAddr string
	flag.StringVar(&listenAddr, "listen", ":1080", "listen addr, e.g: 127.0.0.1:1080")
	flag.Parse()

	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	if err := server.ListenAndServe("tcp", listenAddr); err != nil {
		panic(err)
	}

	logs.Info("listen in %s", listenAddr)
}
