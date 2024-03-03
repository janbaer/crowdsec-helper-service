package main

import (
	"janbaer/crowdsec-helper-service/csclirunner"
	"net"
	"net/http"
	"os"

	"log/slog"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func main() {
	port := 8080
	mux := http.NewServeMux()

	mux.HandleFunc("GET /crowdsec-helper-service/healthcheck", handleHealthCheck)
	mux.HandleFunc("DELETE /crowdsec-helper-service/crowdsec/decisions", handleDeleteCrowdsecDecison)

	logger.Info("Listening on localhost ", "port", port)
	http.ListenAndServe("localhost:8080", mux)
}

func sendResponse(w http.ResponseWriter, statusCode int, body string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, "OK")
}

func handleDeleteCrowdsecDecison(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	ipAddres := queryParams.Get("ip")
	if net.ParseIP(ipAddres) == nil {
		sendResponse(w, http.StatusBadRequest, "The provided IP address is not valid")
		return
	}

	if err := csclirunner.DeleteDecision(ipAddres); err != nil {
		sendResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, "OK")
}
