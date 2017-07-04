package main

import (
	"github.com/antipin1987@gmail.com/rpcj/db"
	"github.com/antipin1987@gmail.com/rpcj/server"
	"net/http"
	"log"
)


func main() {
	defer db.Close()
	go server.StartServerTCP(server.AddrTCP)
	http.HandleFunc("/", server.HttpHandler)
	log.Fatal(http.ListenAndServe(server.AddrHTTP, nil))
}