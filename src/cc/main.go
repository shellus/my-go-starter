package main

import (
	"net/http"
	"github.com/hprose/hprose-golang/util"
	"github.com/pkg/errors"
	"strings"
	"io/ioutil"
	"fmt"
	"sync"
	"os"
	"flag"
	"strconv"
	"time"
)

var u string = ""
var c int = 0
var headersFile string = ""


func main() {
	flag.IntVar(&c, "c", 10, "Concurrent Number")
	flag.StringVar(&u, "u", "https://www.sentris.net/billing/clientarea.php?action=productdetails&id=15327", "url")
	flag.StringVar(&headersFile, "hf", "", "headers file name")

	flag.Parse()

	fmt.Println("concurrentNumber: " + strconv.Itoa(c))
	fmt.Println("url: " + u)
	fmt.Println("headers file name: " + headersFile)
	time.Sleep(time.Second*2)

	for i := 0; i < c; i++ {
		go func() {
			for {
				request()
			}
		}()
	}
	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}

func NewRequest()(req *http.Request, err error){

	req, err = http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	var bytes []byte
	if bytes, err = ioutil.ReadFile(headersFile); err != nil || len(bytes) == 0 {
		fmt.Println("no headers file")
		return nil, err
	}

	lines := strings.Split(string(bytes), "\r\n")
	for _, line := range lines{
		kv := strings.Split(line, ": ")
		req.Header.Set(kv[0], kv[1])
	}
	return
}

func request() {

	var c = http.DefaultClient

	req, err := NewRequest()
	if err != nil {
		fmt.Println(errors.New(err.Error()))
		return
	}

	r, err := c.Do(req)

	if err != nil {
		fmt.Println(errors.New(err.Error()))
		return
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		fmt.Println(errors.New("r.StatusCode is " + util.Itoa(r.StatusCode)))
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)

	html := string(bytes)

	if strings.Index(html, "30/01/2020") == -1 {
		fmt.Println("error")
		fmt.Println("")
		return
	} else {
		fmt.Println("ok")
	}

	if strings.Index(html, "Login to tamil Host failed.") == -1 {
		fmt.Println("fail")
	} else {
		fmt.Println("success !")
		os.Exit(0)
	}
}
