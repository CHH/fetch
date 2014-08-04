package main

import (
	"github.com/CHH/fetch"
	"net/rpc"
	"net/rpc/jsonrpc"
	"net"
	"log"
	"os"
)

func main() {
	fetchService := fetch.NewService(fetch.Config{MaxIdleConnections: 10})
	rpcServer := rpc.NewServer()
	rpcServer.RegisterName("fetch", fetchService)

	os.Remove("/tmp/fetch.sock")

	l, err := net.Listen("unix", "/tmp/fetch.sock")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, _ := l.Accept()
		go rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
