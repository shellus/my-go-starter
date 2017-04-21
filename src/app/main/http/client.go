package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
)

type O struct {
	body_len int
	ir int
	status string
}

func main() {
	i := make(chan int)
	o := make(chan O)
	for ir := 0; ir < 100; ir++ {
		go work(ir, i, o)
	}
	for n:=0; ; n++{
		o := <- o
		fmt.Printf("%.[1]3d: %s %d ) %d\n",o.ir , o.status, o.body_len, n)
	}
}




func work(ir int, i chan int, o chan O) {
	for {
		var url string
		url = "https://laravel-china.org/"
		resp, err := http.Get(url)

		if err != nil {
			fmt.Print(err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Print(err)
			continue
		}

		o <- O{body_len:len(body),ir:ir,status:resp.Status}
	}
}