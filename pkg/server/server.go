package server

import (
	"log"
	"net/http"
	"github.com/Shoetan/pkg/db"

	"github.com/Shoetan/pkg/service"
)




type APISERVER struct{
	addr string
}

func NewAPISERVER(addr string) *APISERVER {
	return &APISERVER{
		addr: addr,
	}
}


func (s *APISERVER) Run() error{

	dele, err := db.Database()
	
	if err != nil {
		log.Fatalf("Cannot connect to the database : %s", err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/register", service.RegisterUser(dele))

	server := http.Server{
		Addr: s.addr,
		Handler: router,
	}

	log.Printf("The server is listening on port %s", s.addr)

	return server.ListenAndServe()
}