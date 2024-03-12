package server

import (
	"fmt"
	"net/http"

	"github.com/NotYourAvgCoder/HTTP-URL-Shortner/handler"
	"github.com/gorilla/mux"
)

type Server struct {
	Port int
}

func (s *Server) Start() error {
	fmt.Println("Creating router !!!")
	router := mux.NewRouter()

	fmt.Println("Creating handler !!!")
	handler := handler.CreateHandler()

	router.HandleFunc("/ping", handler.Ping).Methods(("GET"))
	router.HandleFunc("/create-short-url", handler.CreateShortURL).Methods(("POST"))

	fmt.Printf("Starting server at port : %v\n", s.Port)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", s.Port), router)

	if err != nil {
		return fmt.Errorf("error while trying to start server : %v", err)
	}
	return nil
}
