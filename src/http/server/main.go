package main

import (
	"net/http"
	"encoding/json"
	"os/signal"
	"os"
	"fmt"
	"context"
)

// api 接口返回的通用格式
type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    interface{} `json:"data"`
}

// 登录接口返回的用户基本资料
type UserData struct {
	Name string `json:"name"`
	Id   int `json:"id"`
}

func main() {
	mux := http.NewServeMux()
	serverTls := &http.Server{Addr: ":443", Handler: mux}

	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(`asset`))))

	mux.HandleFunc("/api/login", func(w http.ResponseWriter, req *http.Request) {

		// 定义返回值
		r := ApiResponse{
			"ok",
			"登录成功",
			UserData{
				"shellus",
				1,
			},
		}

		// 序列化成json
		j, e := json.Marshal(r)
		if e != nil {
			panic(e)
		}

		// 发送数据
		w.Write(j)
	})

	go func() {
		err := serverTls.ListenAndServeTLS("letsencrypt/local/cert1.pem", "letsencrypt/local/privkey1.pem")
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()


	fmt.Printf("serverTls Listen: %s!\n", serverTls.Addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	err := serverTls.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}


}

