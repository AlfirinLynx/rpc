package server

import (
	"net/rpc"
	"net"
	"log"
	"net/rpc/jsonrpc"
	"github.com/antipin1987@gmail.com/rpcj/models"
)

func StartServerTCP(addr string) {
	server := rpc.NewServer()
	for name, model := range models.Registry() {
		server.RegisterName(name, model)
	}
	server.HandleHTTP("/rpc", rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			log.Printf("new connection established\n")
			//io.Copy(os.Stdout, conn)
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}