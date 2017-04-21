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
		url = "https://segmentfault.com/questions"
		resp, err := http.Get(url)

		if err != nil {
			fmt.Print(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)

		o <- O{body_len:len(body),ir:ir,status:resp.Status}
	}
}