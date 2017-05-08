package proxyPool

import (
	"github.com/hprose/hprose-golang/util"
	"net/url"
	"io/ioutil"
	"strings"
	"net/http"
	"errors"
)

// 验证代理ip有效性及匿名性
func VerifyProxy(addr string) (result bool, err error) {
	uuid := util.UUIDv4()

	u, err := url.Parse("http://" + addr);
	if err != nil {
		return
	}
	c := &http.Client{Transport:&http.Transport{Proxy:http.ProxyURL(u)}}

	r, err := c.Get("http://tz.endaosi.com/?uuid=" + uuid)

	if err != nil {
		return
	}
	if r.StatusCode != 200 {
		err = errors.New("HTTP status code is no 200")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	if len(body) == 0 {
		err = errors.New("Body length is zero")
		return
	}
	if strings.Index(string(body), uuid) == -1 {
		err = errors.New("Body not found uuid")
		return
	}
	result = true
	return
}