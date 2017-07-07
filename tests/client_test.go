package tests

import (
	"net"
	"net/rpc/jsonrpc"
	"testing"
	"fmt"
	"time"
	"github.com/antipin1987@gmail.com/rpcj/config"
	"github.com/antipin1987@gmail.com/rpcj/server"
	"github.com/antipin1987@gmail.com/rpcj/db"
	"github.com/antipin1987@gmail.com/rpcj/models/orm"
	"github.com/stretchr/testify/assert"
	"net/rpc"
	"log"
	"github.com/antipin1987@gmail.com/rpcj/models"
)
var addr string
var c *rpc.Client
var rpcServer *rpc.Server

func init() {
	rpcServer = rpc.NewServer()
	conf := config.Get().Sub("net.tcp")
	addr = fmt.Sprintf("%s:%s", conf.GetString("host"), conf.GetString("port"))
	go server.StartServerTCP(rpcServer, addr)
	time.Sleep(3 * time.Second)

	client, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	c = jsonrpc.NewClient(client)
}

func Test(t *testing.T) {
	var rep bool
	arg := "kvAnt"
	err := c.Call("User.Add", arg, &rep)
	if err != nil {
		t.Error("User.Add error:", err)
		return
	}
	fmt.Println("Success: ", rep)
	assert.True(t, rep)
	var count int
		if err := db.DB().Model(&orm.User{}).Where("login = ?", arg).Count(&count).Error; err != nil {
			t.Error(err)
		}
	fmt.Println("Records in db: ", count)
	assert.True(t, count > 0)
}


func TestUserFind(t *testing.T) {
	f := &models.Filter{Login: "Cat", Date: time.Now()}
	var usrs = make([]orm.User, 0)


	err := c.Call("User.Find", f, &usrs)
	if err != nil {
		t.Error("User.Find error: ", err)
		return
	}
	fmt.Println("Success: ", usrs)

}