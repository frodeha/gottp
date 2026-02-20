package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/frodeha/gottp/http"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	// defer lis.Close()

	server := http.NewServer(lis)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGABRT)
	go func() {
		<-stop
		err := server.Close()
		if err != nil {
			panic(err)
		}
	}()

	err = server.ServeHTTP()
	if err != nil {
		panic(err)
	}
}
