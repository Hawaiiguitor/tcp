package main

import (
	"fmt"

	"github.com/Hawaiiguitor/tcp/tcpClient/client"
)

func main() {
	c := client.NewClient()
	err := c.Connect("localhost", "9191")
	if err != nil {
		return
	}
	err = c.Sendfile("./test")
	if err != nil {
		fmt.Print("tcp client failed to send file\n")
	}
	c.Close()
}
