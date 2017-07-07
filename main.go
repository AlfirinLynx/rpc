package main

import (
	"github.com/antipin1987@gmail.com/rpcj/db"
	"github.com/antipin1987@gmail.com/rpcj/server"
	"net/http"
	"log"
	"net/rpc"
)


func main() {
	defer db.Close()

	rpcServer := rpc.NewServer()
	go server.StartServerTCP(rpcServer, server.AddrTCP)
	http.HandleFunc("/", (&server.Serv{Server: rpcServer}).HttpHandler)
	log.Fatal(http.ListenAndServe(server.AddrHTTP, nil))
}