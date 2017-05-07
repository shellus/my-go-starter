package proxyPool

import (
	"net/url"
	"proxyPool/dataSource/cn_proxy_com"
	"github.com/astaxie/beego/logs"
)


var c chan string
func Init(){
	cn_proxy_com.PullData()
	c = cn_proxy_com.Factory()
}
func GetProxyUrl() *url.URL{
	s := <- c
	proxyUrl, _ := url.Parse("http://" + s)
	logs.Info("use proxy addr: " + s)
	return proxyUrl
}