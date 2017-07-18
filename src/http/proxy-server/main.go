package proxy_server

import (
	"net/http"
	"github.com/shellus/pkg/logs"
	"net/http/httputil"
)

func Main(){
	handleFun := func(w http.ResponseWriter,r *http.Request){
		body, err := httputil.DumpRequest(r, true)
		if err != nil {
			logs.Fatal(err)
		}

		r.AddCookie(&http.Cookie{Name:"a", Value:"1"})
		w.Write(body)
	}
	http.HandleFunc("/", handleFun)


	//err := http.ListenAndServe(":80", nil)
	err := http.ListenAndServeTLS(":443", "cert/fullchain.pem", "cert/privkey.pem", nil)
	if err != nil {
		logs.Fatal(err)
	}
}