package server

import (
	"net"
	"fmt"
)

type server struct {

}

func NewServer() *server {
	return &server{}
}

func (srv *server) ListenAndServer(port string) error {
    addr := "0.0.0.0" + ":" + port
	ln, err := net.Listen("tcp", addr)
	if err != nil {
        fmt.Printf("Fail to listening, err: %v\n", err)
		return err
	}
	fmt.Printf("server listening on: %s\n", addr)
    for {
		rw, err := ln.Accept()
		if err != nil {
			fmt.Printf("Fail to accept incoming conn, err: %v\n", err)
			continue
		}
		go handerConn(rw)
	}
}

func handerConn(rw net.Conn) {
    defer rw.Close()
}