package main
import (

	rpc "github.com/hprose/hprose-golang/rpc/fasthttp"
	"github.com/valyala/fasthttp"
	"fmt"
)

func hello(name string) string {
	return "Hello !"
}

func main() {
	service := rpc.NewFastHTTPService()
	service.AddFunction("hello", hello)
	fmt.Println("server started.")
	fasthttp.ListenAndServe(":6060", service.ServeFastHTTP)

}