package main

import (
	"net/http"
	"fmt"
	"net/http/httptrace"
)

func main() {

	req, err := http.NewRequest("GET", "http://mirrors.163.com/centos/6.9/isos/x86_64/CentOS-6.9-x86_64-bin-DVD1.iso", nil)
	if err != nil {
		panic(err)
	}

	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	resp, err := http.DefaultClient.Do(req)

	fmt.Println(resp.Header)
	//body, err := ioutil.ReadAll(resp.Body)

	//if err != nil {
	//	fmt.Print(err)
	//	continue
	//}
}
