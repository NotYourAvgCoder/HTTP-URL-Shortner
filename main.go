package main

import (
	"fmt"

	"github.com/NotYourAvgCoder/HTTP-URL-Shortner/server"
)

func main() {
	fmt.Println("Starting !!!")
	serv := server.Server{Port: 3030}
	serv.Start()
}
