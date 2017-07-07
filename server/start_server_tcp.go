package server

import (
	"net/rpc"
	"net"
	"log"
	"net/rpc/jsonrpc"
	"github.com/antipin1987@gmail.com/rpcj/models"
)

func StartServerTCP(rpcServer *rpc.Server, addr string) {
	for name, model := range models.Registry() {
		rpcServer.RegisterName(name, model)
	}

	listener, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			log.Printf("new connection established\n")
			go rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}