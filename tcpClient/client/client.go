package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
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

func (c *client) Send(data string) {
	_, err := c.conn.Write([]byte(data))
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
