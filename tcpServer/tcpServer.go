package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Hawaiiguitor/tcp/tcpServer/server"
)

func main() {
	srv := server.NewServer()
	err := srv.ListenAndServer("9191")
	if err != nil {
		panic(err)
	}

	var captureSignal = make(chan os.Signal, 1)
	signal.Notify(captureSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	<-captureSignal
	srv.Stop()
}
