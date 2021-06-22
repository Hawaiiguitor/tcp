package main

import (
	"time"

	"github.com/Hawaiiguitor/tcp/tcpClient/client"
)

func main() {
	c := client.NewClient()
	c.Connect("localhost", "9191")
	c.Send("test\n")

	for {
		time.Sleep(1 * time.Second)
	}
}
