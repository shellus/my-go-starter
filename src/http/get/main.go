package main

import (
	"net/http"
	"proxyPool"
	"io/ioutil"
	"log"
)

func main() {
	proxyPool.Init()

	c := &http.Client{Transport:&http.Transport{Proxy:http.ProxyURL(proxyPool.GetProxyUrl())}}

	r, err := c.Get("http://ip.taobao.com/service/getIpInfo.php?ip=myip")
	if err!= nil {
		log.Fatalln(err)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err!= nil {
		log.Fatalln(err)
		return
	}
	log.Println(string(b))
}
