package main

import (
	"net/http"
	"fmt"
)

func main() {
	http.Handle("/asset", http.FileServer(http.Dir(`C:\data\c\frp`)))

	http.Handle("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "哈哈哈")
	})

	http.ListenAndServe(":8080", nil)
}

