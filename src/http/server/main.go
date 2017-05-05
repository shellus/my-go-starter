package main

import (
	"net/http"
	"encoding/json"
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
	http.Handle("/asset/", http.StripPrefix("/asset/", http.FileServer(http.Dir(`asset`))))

	http.HandleFunc("/api/login", func(w http.ResponseWriter, req *http.Request) {

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

	http.ListenAndServe(":8080", nil)
}

