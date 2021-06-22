package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
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
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Fail to accept incoming conn, err: %v\n", err)
			continue
		}
		go handerConn(conn)
	}
}

func handerConn(conn net.Conn) {
	defer conn.Close()
	rd := bufio.NewReader(conn)
	wd := bufio.NewWriter(conn)
	rw := bufio.NewReadWriter(rd, wd)
	for {
		data, _, err := rw.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("Fail to read data, err: %v", err)
		}
		fmt.Printf("data: %s\n", string(data))
		sendData := append(data, '\n')

		_, err = conn.Write(sendData)
		if err != nil {
			fmt.Printf("Send data failure, err: %v", err)
			return
		}
		break
	}
}
