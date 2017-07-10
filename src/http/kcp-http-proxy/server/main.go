package main

import (
	"net/http"
	"log"
	"github.com/elazarl/goproxy"
	"github.com/xtaci/kcp-go"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			return r, nil
		})

	http.Handle("/", proxy)

	server := &http.Server{Addr: ":8080", Handler: nil}
	err := ListenAndServe(server)

	if err != nil {
		log.Fatal(err)
	}
}
func ListenAndServe(srv *http.Server) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := kcp.ListenWithOptions(addr, nil, 10, 3)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}