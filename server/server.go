package server

import (
	"github.com/gsq/music_bakcend_micorservice/myredis"
	"log"
	"net/http"
)

type Server struct {
	ip      string
	port    string
	handler http.Handler
}

func NewServer(ip, port string, handler http.Handler) *Server {
	return &Server{
		ip:      ip,
		port:    port,
		handler: handler,
	}
}

func (s *Server) Run() error {

	httpServ := &http.Server{
		Addr:    s.ip + ":" + s.port,
		Handler: s.handler,
	}

	myredis.InitRedis()

	log.Printf("Server start and listen at %s.", httpServ.Addr)
	if err := httpServ.ListenAndServe(); err != nil {
		log.Fatalf("Server start fail: %s", err)
	}

	return nil
}
