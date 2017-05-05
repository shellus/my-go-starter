package cn_proxy_com

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"strconv"
)

type ProxyAddr struct {
	Ip   string
	Port int
}

var s = []ProxyAddr{}

func Factory() (o chan ProxyAddr) {
	o = make(chan ProxyAddr)
	go func() {
		for {
			for _, n := range s {
				o <- n
			}
		}
	}()
	return
}

func PullData() (err error) {

	res, err := http.Get("http://cn-proxy.com/")
	if err != nil {
		return
	}


	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}


	r := regexp.MustCompile(`<td>((\d{1,3}\.){3}\d{1,3})</td>\n<td>(\d{1,6})</td>`).FindAllSubmatch(html, -1)
	for _, l := range r {


		p, err := strconv.Atoi(string(l[3]))
		if err != nil {
			return
		}



		s = append(s, ProxyAddr{Ip:string(l[1]), Port:p})
	}
	return
}


