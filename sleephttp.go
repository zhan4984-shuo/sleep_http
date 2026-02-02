package main

import (
	"encoding/json" // CHANGED
	"fmt"
	"log"
	"mime" // CHANGED
	"net/http"
	"os"
	"strconv"
    "time"
)

type RequestBody struct { // CHANGED
	Name string `json:"name"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("helloworld: received a %s request", r.Method) // CHANGED

	// CHANGED: 只接受 POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed: use POST", http.StatusMethodNotAllowed)
		return
	}

	// CHANGED: 只接受 application/json（允许带 charset）
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

	// CHANGED: 解析 JSON body
	var body RequestBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // CHANGED: 多余字段直接报错
	if err := dec.Decode(&body); err != nil {
		http.Error(w, "Bad Request: invalid JSON body", http.StatusBadRequest)
		return
	}

	// CHANGED: 确保 body 里只有一个 JSON 对象（防止 `{} {}` 这种）
	if dec.More() {
		http.Error(w, "Bad Request: multiple JSON values", http.StatusBadRequest)
		return
	}

	// CHANGED: 参数校验
	if body.Name == "" {
		http.Error(w, "Bad Request: missing 'name'", http.StatusBadRequest)
		return
	}

	i, err := strconv.ParseInt(body.Name, 10, 64) 
	
	if err != nil {
		fmt.Printf("Error during conversion: %v\n", err)
		return
	}

	fmt.Printf("Sleeping for %d milliseconds...\n", i)

	time.Sleep(time.Duration(i) * time.Millisecond)
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