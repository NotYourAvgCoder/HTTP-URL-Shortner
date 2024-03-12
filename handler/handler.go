package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Handler struct{}

func (h *Handler) Ping(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Write([]byte("Hello world !!!"))
}

func (h *Handler) CreateShortURL(response http.ResponseWriter, request *http.Request) {
	url, err := getValueFromBody("url", request.Body)

	if err != nil {
		http.Error(response, fmt.Sprintf("error while trying to fetch url value from request body : %v", err), http.StatusBadRequest)
	}

	fmt.Printf("URL value retrevived : %v\n", url)

	response.WriteHeader(200)
	response.Write([]byte(url))
}

func getValueFromBody(key string, reqBody io.ReadCloser) (string, error) {
	body, err := io.ReadAll(reqBody)

	if err != nil {
		return "", err
	}

	// Parse request body as JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	// Retrieve the value of the "url" key
	val, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("missing key %v from request body", key)
	}

	return val, nil

}
