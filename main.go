package main

import (
	"github.com/antipin1987@gmail.com/rpcj/db"
	"github.com/antipin1987@gmail.com/rpcj/server"
	"net/http"
	"log"
	"net/rpc"
)


//json-rpc сервер, слушающий на двух портах: один - на чистом tcp, другой - по http
//порты прописаны в конфиге (см. файл config.yaml)

func main() {
	defer db.Close()

	rpcServer := rpc.NewServer()
	//Старт rpc сервера на чистом tcp(не обрабатывает http-запросы)
	go server.StartServerTCP(rpcServer, server.AddrTCP)

	//И одновременно старт сервера, слушающего http-запросы(на другом порте) с помощью структуры-адаптера
	http.HandleFunc("/", (&server.Serv{Server: rpcServer}).HttpHandler)
	log.Fatal(http.ListenAndServe(server.AddrHTTP, nil))
}