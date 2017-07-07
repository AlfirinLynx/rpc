package server

import (
	"net/http"
	"fmt"
	"github.com/antipin1987@gmail.com/rpcj/config"
	"log"
	"io"
	"net/rpc/jsonrpc"
	"net/rpc"
)

var AddrTCP, AddrHTTP string

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

type Serv struct {
	Server *rpc.Server
}

func (s *Serv) HttpHandler(w http.ResponseWriter, r *http.Request) {
	serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	err := s.Server.ServeRequest(serverCodec)
	if err != nil {
		log.Printf("Error while serving JSON request: %v", err)
		http.Error(w, "Error while serving JSON request, details have been logged.", 500)
		return
	}
}







func Copy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}





func init() {
	conf := config.Get().Sub("net.tcp")
	AddrTCP = fmt.Sprintf("%s:%s", conf.GetString("host"), conf.GetString("port"))
	conf = config.Get().Sub("net.http")
	AddrHTTP = fmt.Sprintf("%s:%s", conf.GetString("host"), conf.GetString("port"))
}