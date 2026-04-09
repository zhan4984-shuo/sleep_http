package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"time"
)

type RequestBody struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	//log.Printf("helloworld: received a %s request", r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed: use POST", http.StatusMethodNotAllowed)
		return
	}

	mediaType := r.Header.Get("Content-Type")
	if mediaType == "" {
		http.Error(w, "Unsupported Media Type: Content-Type required", http.StatusUnsupportedMediaType)
		return
	}
	mt, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		http.Error(w, "Bad Request: invalid Content-Type", http.StatusBadRequest)
		return
	}
	if mt != "application/json" {
		http.Error(w, "Unsupported Media Type: use application/json", http.StatusUnsupportedMediaType)
		return
	}

	var body RequestBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		http.Error(w, "Bad Request: invalid JSON body", http.StatusBadRequest)
		return
	}

	//log.Printf("uuid: %s", body.UUID)
	//log.Printf("uuid: %s, x-request-id: %s", body.UUID, r.Header.Get("X-Request-ID"))

	if dec.More() {
		http.Error(w, "Bad Request: multiple JSON values", http.StatusBadRequest)
		return
	}

	if body.Name == "" {
		http.Error(w, "Bad Request: missing 'name'", http.StatusBadRequest)
		return
	}

	i, err := strconv.ParseInt(body.Name, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request: invalid 'name'", http.StatusBadRequest)
		return
	}

	//log.Printf("Sleeping for %d milliseconds...\n", i)
	time.Sleep(time.Duration(i) * time.Millisecond)

	// 把 request body 原样返回
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func main() {
	log.Print("helloworld: starting server...")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
