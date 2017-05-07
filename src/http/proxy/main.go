package main

import (
	"net/http"
	"log"
	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			host := "www.google.com"

			log.Println(r.Header.Get("Host"))

			r.Header.Set("Host", host)
			r.Host = host
			r.URL.Host = host
			r.URL.Scheme = "https"
			return r, nil
		})

	http.Handle("/", proxy)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
