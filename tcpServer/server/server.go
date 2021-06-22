package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/Hawaiiguitor/tcp/message"
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
		hd, _, err := rw.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("Fail to read data, err: %v", err)
		}

		header, err := message.DecodeHeader(hd)
		datasize := header.DataSize
		if datasize > message.MAX_NET_DATA_SIZE {
			rw.Write([]byte("Deny"))
			break
		}
		if header.OpCode == message.OP_SENDINFO {
			buf := make([]byte, datasize)
			_, err := io.ReadFull(rw, buf)
			if err != nil {
				log.Fatalf("SendInfo: Fail to read body, %v", err)
			}
			fname := string(buf)
			fd, err := os.Create(fname)
			if err != nil {
				log.Fatal("Fail to create file")
			}
			fmt.Printf("Create file: %s\n", fd.Name())
			defer fd.Close()

			_, err = conn.Write([]byte("Access"))
			if err != nil {
				fmt.Printf("Send data failure, err: %v", err)
				return
			}
		}

		break
	}
}
