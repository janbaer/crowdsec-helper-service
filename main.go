package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", handleHealthCheck)
	mux.HandleFunc("DELETE /crowdsec/decisions", handleDeleteCrowdsecDecison)

	fmt.Println("Listening on port 8080")
	http.ListenAndServe("localhost:8080", mux)
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// You can also write a response body if needed
	w.Write([]byte("This is a response with status code 200"))
}

func handleDeleteCrowdsecDecison(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
