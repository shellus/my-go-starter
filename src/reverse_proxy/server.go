package main

import (
	"net/http/httputil"
	"net/http"
	"net/url"
	"strings"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery

	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

		req.Host = target.Host
	}

	return &httputil.ReverseProxy{Director: director}
}

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

// 删除cookie域名限制
func deleteCookiesDomain(resp *http.Response){
	var cookies []*http.Cookie
	cookies = []*http.Cookie{}
	for _, c := range resp.Cookies(){
		c.Domain = ""
		cookies = append(cookies, c)
	}
	resp.Header.Del("Set-Cookie")
	for _, c := range cookies{
		resp.Header.Add("Set-Cookie", c.String())
	}
}

// 修正跳转
func alterationLocation(resp *http.Response){
	if r, err := resp.Location(); err == nil {
		r.Host = "127.0.0.1:8001"
		resp.Header.Set("Location", r.String())
	}
}

var proxy = func(_ *http.Request) (*url.URL, error) {
	return url.Parse("http://127.0.0.1:1080")
}

var proxyTransport = http.Transport{Proxy:proxy}

func main() {
	backendUrl, err := url.Parse("https://segmentfault.com/")
	if err != nil {
		panic(err)
	}
	proxyHandle := NewSingleHostReverseProxy(backendUrl)
	proxyHandle.Transport = RoundTripperFunc(func(req *http.Request) (*http.Response, error) {

		// 匿名代理
		req.Header.Del("X-Forwarded-For")

		resp, err := proxyTransport.RoundTrip(req)
		if err != nil {
			return resp, err
		}

		alterationLocation(resp)

		deleteCookiesDomain(resp)

		return resp, nil
	})
	err = http.ListenAndServe(":8001", proxyHandle)
	if err != nil {
		panic(err)
	}
}