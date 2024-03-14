package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"

	inmemorystorage "github.com/NotYourAvgCoder/HTTP-URL-Shortner/in_memory_storage"
	"github.com/gorilla/mux"
)

type MutexCounter struct {
	m   sync.Mutex
	val uint
}

type Handler struct {
	counter  *MutexCounter
	database *inmemorystorage.RedisDatabase
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

func CreateHandler() *Handler {
	repo := inmemorystorage.InitializeRedisDB("0.0.0.0:6379", "", 0)
	err := repo.Connect()
	if err != nil {
		log.Fatal().Msgf("error while connecting to redis server : %v", err)
	}
	return &Handler{
		counter:  &MutexCounter{},
		database: repo,
	}
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

	shortenedURL := fmt.Sprintf("http://localhost:3030/url/%v", h.counter.Get())

	err = h.database.Insert(shortenedURL, url)

	if err != nil {
		http.Error(response, fmt.Sprintf("error while trying to store url : %v", err), http.StatusInternalServerError)
	}

	response.WriteHeader(200)
	response.Write([]byte(shortenedURL))
}

func (h *Handler) RedirectTo(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["id"]

	shortenedURL := fmt.Sprintf("http://localhost:3030/url/%v", key)

	log.Info().Msgf("Shortened URL received key : %v", key)

	val, err := h.database.Get(shortenedURL)

	log.Info().Msgf("Shortened URL received value : %v", val)

	if err != nil {
		http.Error(response, fmt.Sprintf("error while trying to fetch url : %v", err), http.StatusInternalServerError)
	}
	http.Redirect(response, request, val, http.StatusAccepted)
}

func getValueFromBody(key string, reqBody io.ReadCloser) (string, error) {
	// So here we're reading byte stream of request body
	body, err := io.ReadAll(reqBody)

	if err != nil {
		return "", err
	}

	// Parse request body as JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	// Retrieve the value of the "url" key and explicity convert it to string
	val, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("missing key %v from request body", key)
	}

	return val, nil

}
