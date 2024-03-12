package handler

import "net/http"

type Handler struct{}

func (h *Handler) Ping(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Write([]byte("Hello world !!!"))
}
