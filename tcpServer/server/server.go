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
	ln   net.Listener
	quit chan interface{}
}

func NewServer() *server {
	return &server{}
}

func (srv *server) Stop() {
	srv.ln.Close()
}

func (srv *server) ListenAndServer(port string) error {
	addr := "0.0.0.0" + ":" + port
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Fail to listening, err: %v\n", err)
		return err
	}
	srv.ln = ln
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
		rw.Write([]byte("Deny\n"))
		return
	}

	var fname string
	if header.OpCode == message.OP_SENDINFO {
		buf := make([]byte, datasize)
		_, err := io.ReadFull(rw, buf)
		if err != nil {
			log.Fatalf("SendInfo: Fail to read body, %v", err)
		}
		fname = string(buf)
		fd, err := os.Create(fname)
		if err != nil {
			log.Fatal("Fail to create file")
		}
		fmt.Printf("Create file: %s\n", fd.Name())

		_, err = conn.Write([]byte("Access\n"))
		if err != nil {
			fmt.Printf("Send data failure, err: %v", err)
			return
		}
		fd.Close()
	} else {
		fmt.Printf("Missing file info from client, interrupt this connection\n")
		_, err = conn.Write([]byte("Deny\n"))
		if err != nil {
			fmt.Printf("Send data failure, err: %v", err)
		}
		return
	}

	buf := make([]byte, 1024)
	for datasize > 0 {
		n, err := rw.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Printf("Fail to read data, err: %v, n: %d", err, n)
			return
		}
		fmt.Printf("len: %d, fname: %s \n", len(buf[0:n]), fname)
		f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			fmt.Printf("Fail to open file, err: %v", err)
			return
		}
		_, err = f.Write(buf[0 : n+1])
		if err != nil {
			fmt.Printf("Fail to write data, err: %v", err)
			return
		}
		datasize -= n
	}

}
