package main

import (
	"github.com/CHH/fetch"
	"net/rpc"
	"net/rpc/jsonrpc"
	"net"
	"log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fetchService := fetch.NewService(fetch.Config{MaxIdleConnections: 10})
	rpcServer := rpc.NewServer()
	rpcServer.RegisterName("fetch", fetchService)

	l, err := net.Listen("tcp", ":1337")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, _ := l.Accept()
		go rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
