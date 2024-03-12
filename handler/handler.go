package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type MutexCounter struct {
	m   sync.Mutex
	val uint
}

type Handler struct {
	counter *MutexCounter
}

func (mc *MutexCounter) Inc() uint {
	mc.m.Lock()
	defer mc.m.Unlock()
	mc.val += 1
	return mc.val
}

func (mc *MutexCounter) Get() uint {
	mc.m.Lock()
	defer mc.m.Unlock()
	return mc.val
}

func (h *Handler) Ping(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Write([]byte("Hello world !!!"))
}

/*
*
  - curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"url":"www.google.com"}' \
    http://localhost:3030/create-short-url

*
*/
func (h *Handler) CreateShortURL(response http.ResponseWriter, request *http.Request) {
	url, err := getValueFromBody("url", request.Body)

	if err != nil {
		http.Error(response, fmt.Sprintf("error while trying to fetch url value from request body : %v", err), http.StatusBadRequest)
	}

	fmt.Printf("URL value retrevived : %v\n", url)
	h.counter.Inc()

	response.WriteHeader(200)
	response.Write([]byte(fmt.Sprintf("http://localhost:3030/%v", h.counter.Get())))
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
