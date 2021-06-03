package main

import (
	""
)

func main(){
	srv := server.NewServer()
	err := srv.ListenAndServer("9191")
	if err != nil {
		panic(err)
	}
}