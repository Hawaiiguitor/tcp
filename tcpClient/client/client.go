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
}

func NewClient() client {
	return client{}
}

func (c *client) Connect(ip, port string) {
	addr := ip + ":" + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Fail to connect %s, err: %v", addr, err)
	}
	c.conn = conn
}

func (c *client) Sendfile(fname string) {
	fd, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Fail to open file %s", fname)
	}
	f_info, err := fd.Stat()
	if err != nil {
		log.Fatalf("Fail to get fileinfo of  %s", fname)
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
	}

	_, err = c.conn.Write([]byte(data))
	if err != nil {
		fmt.Printf("Fail to send data, err: %v", err)
		return
	}

	rd := bufio.NewReader(c.conn)
	recv, _, err := rd.ReadLine()
	if err != nil {
		if err == io.EOF {
			fmt.Printf("recv: %s", string(recv))
			return
		}
		fmt.Printf("Fail to receive data, err: %v", err)
	}
	fmt.Printf("recv: %s", string(recv))
}

func (c *client) Close() {
	c.conn.Close()
}
