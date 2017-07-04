package tests

import (
	"net"
	"net/rpc/jsonrpc"
	"testing"
	"fmt"
	"time"
	"github.com/antipin1987@gmail.com/rpcj/config"
	"github.com/antipin1987@gmail.com/rpcj/server"
)


func Test(t *testing.T) {
	conf := config.Get().Sub("net.tcp")
	addr := fmt.Sprintf("%s:%s", conf.GetString("host"), conf.GetString("port"))
	go server.StartServerTCP(addr)
	var rep bool
	time.Sleep(3 * time.Second)
	client, err := net.Dial("tcp", addr)
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


