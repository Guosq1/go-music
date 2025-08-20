package main

import (
	"github.com/gsq/music_bakcend_micorservice/pkg/apiGateway"
	"github.com/gsq/music_bakcend_micorservice/server"
)

func main() {

	ip := "127.0.0.1"
	port := "8080"

	engine := apiGateway.Setup()

	srv := server.NewServer(ip, port, engine)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
