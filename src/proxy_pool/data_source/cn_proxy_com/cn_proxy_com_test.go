package cn_proxy_com_test

import (
	"testing"
	"proxy_pool/data_source/cn_proxy_com"
	"proxy_pool"
)

func TestGet(t *testing.T) {
	addrs, err := cn_proxy_com.Get()
	if err != nil {
		t.Error(err)
	}
	for _, addr := range addrs {
		b, err := proxy_pool.VerifyProxy(addr);
		if b == false {
			t.Error(err)
		}
	}
}


