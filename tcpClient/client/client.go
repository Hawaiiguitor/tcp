package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/Hawaiiguitor/tcp/message"
)

type client struct {
	conn net.Conn
	Fin  chan struct{}
}

func NewClient() client {
	return client{
		Fin: make(chan struct{}),
	}
}

func (c *client) Connect(ip, port string) error {
	addr := ip + ":" + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Fail to connect %s, err: %v", addr, err)
		return err
	}
	c.conn = conn
	return nil
}

func (c *client) Sendfile(fname string) error {
	fd, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Fail to open file %s", fname)
		return err
	}
	f_info, err := fd.Stat()
	if err != nil {
		log.Fatalf("Fail to get fileinfo of  %s", fname)
		return err
	}

	// Send fileinfo
	msg := &message.TcpMsg{}
	msg.Header = message.MsgHeader{
		DataSize: len(f_info.Name()),
		OpCode:   message.OP_SENDINFO,
	}
	msg.Body = []byte(f_info.Name())
	data, err := message.ConstructMsg(msg)
	if err != nil {
		log.Fatalf("Fail to construct message")
		return err
	}

	_, err = c.conn.Write([]byte(data))
	if err != nil {
		fmt.Printf("Fail to send fileinfo, err: %v", err)
		return err
	}

	rd := bufio.NewReader(c.conn)
	recv, _, err := rd.ReadLine()
	if err != nil {
		if err == io.EOF {
			fmt.Printf("recv: %s", string(recv))
			return err
		}
		fmt.Printf("Fail to receive data from server, err: %v", err)
	}
	fmt.Printf("recv: %s from server\n", string(recv))
	switch string(recv) {
	case "Access":
		// send file
		fmt.Println("start to send file content")
		sendFileContent(fd, int(f_info.Size()), c.conn)
	default:
		fmt.Print("Receive unknown result")
	}

	return nil
}

func (c *client) Close() {
	c.conn.Close()
}

func sendFileContent(fd *os.File, size int, conn net.Conn) error {

	buf := make([]byte, 1024)
	for size > 0 {
		n, err := fd.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Print("Fail to read data, err", err)
			return err
		}
		_, err = conn.Write([]byte(buf[0 : n+1]))
		if err != nil {
			fmt.Printf("Fail to send fileinfo, err: %v", err)
			return err
		}
		size -= n
	}
	return nil
}
