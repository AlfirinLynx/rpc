package server

import (
	"net/http"
	"net"
	"fmt"
	"github.com/antipin1987@gmail.com/rpcj/config"
	"log"
	"io"
)

var AddrTCP, AddrHTTP string

func HttpHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := net.Dial("tcp", AddrTCP)
	if err != nil {
		log.Fatal(err)
	}
	Copy(conn, r.Body)
	fmt.Println("Sent body")
	buf := make([]byte, 4096)
	conn.Read(buf)
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
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