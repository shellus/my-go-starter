package main

import (
	"net/http"
	"github.com/shellus/pkg/logs"
	"flag"
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
	"net"
	"time"
	"bufio"
	"strings"
	"os"
)

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

type StaticCredentials map[string]string

func parseCredFile(filename string) StaticCredentials{
	var creds = make(StaticCredentials)
	f, err := os.Open(filename)
	if err != nil {
		logs.Fatal("open cred_file err: %s", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), "=")

		if len(arr) != 1 && len(arr) != 2 {
			logs.Info("Ignore Credential: %s", scanner.Text())
			continue
		}

		var (
			user string = ""
			passwd string = ""
		)

		user = arr[0]
		if len(arr) == 2 {
			passwd = arr[1]
		}

		creds[user] = passwd
	}
	if err := scanner.Err(); err != nil {
		logs.Fatal("scan cred_file err: %s", err)
	}
	return creds
}

var (
	listenAddr string
	credFile string
	creds StaticCredentials
)

func main() {
	loadConfig()
	startHttpProxy()
}

func loadConfig(){
	flag.StringVar(&listenAddr, "listen", ":8080", "listen addr, e.g: 127.0.0.1:8080")
	flag.StringVar(&credFile, "cred_file", "./cred.ini", "a INI File, Content e.g: user1=foo\nuser2=bar")
	flag.Parse()

	creds = parseCredFile(credFile)
}
func authMethod(user, passwd string)bool{
	logs.Debug("Attempt Authenticate user: %s passwd: %s", user, passwd)
	return user != "" && creds[user] == passwd
}
func startHttpProxy(){
	httpProxy := goproxy.NewProxyHttpServer()

	httpProxy.Verbose = true

	httpProxy.OnRequest().Do(auth.Basic("my_realm", authMethod))

	http.Handle("/", httpProxy)

	server := &http.Server{Addr: listenAddr, Handler: nil}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		logs.Fatal("Listen err: %s", err)
	}

	logs.Info("HTTP Proxy Server Listen In %s", listenAddr)

	err = server.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
	if err != nil {
		logs.Fatal(err)
	}
}