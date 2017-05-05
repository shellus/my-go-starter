package proxyPool

import (
	"proxyPool/dataSource/cn_proxy_com"
)

type ProxyAddr cn_proxy_com.ProxyAddr

var c ProxyAddr

func Init(){
	cn_proxy_com.PullData()
	c = cn_proxy_com.Factory()
}
func GetOne() ProxyAddr{
	return <- c
}
