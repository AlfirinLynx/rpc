package tests

import (
	"net/rpc"
	"net"
	"log"
	"net/rpc/jsonrpc"
	"github.com/antipin1987@gmail.com/rpcj/db"
	"github.com/antipin1987@gmail.com/rpcj/models"
	"testing"
	"fmt"
	"time"
)


func Test(t *testing.T) {
	go startServer()
	var rep bool
	time.Sleep(3 * time.Second)
	client, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		t.Error(err)
	}

	c := jsonrpc.NewClient(client)
	err = c.Call("User.Add", "kvAnt", &rep)
	if err != nil {
		t.Error("User.Add error:", err)
		return
	}
	fmt.Println("Success: ", rep)
}

func startServer() {
	defer db.Close()
	server := rpc.NewServer()
	for name, model := range models.Registry() {
		server.RegisterName(name, model)
	}
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp", ":3000")
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
