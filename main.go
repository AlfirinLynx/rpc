package main

import (
	"github.com/antipin1987@gmail.com/rpcj/db"
	"github.com/antipin1987@gmail.com/rpcj/models"
	"net/rpc"
	"net/rpc/jsonrpc"
	"net"
	"log"
)

func main() {
	defer db.Close()
	server := rpc.NewServer()
	for name, model := range models.Registry() {
		server.RegisterName(name, model)
	}
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			log.Printf("new connection established\n")
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}