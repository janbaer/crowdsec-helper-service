package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"janbaer/crowdsec-helper-service/csclirunner"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"log/slog"
)

var logger *slog.Logger

var (
	version = "dev" // this variable holds the current build version
)

func init() {
	// logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func main() {
	var port int

	flag.IntVar(&port, "p", 8000, "Provide a port number")
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /crowdsec-helper-service/healthcheck", handleHealthCheck)
	mux.HandleFunc("DELETE /crowdsec-helper-service/crowdsec/decisions", handleDeleteCrowdsecDecison)

	logger.Info("Listening on localhost ", "port", port)
	http.ListenAndServe(fmt.Sprintf("localhost:%d", port), mux)
}

func sendResponse(w http.ResponseWriter, statusCode int, body string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	healthData := map[string]string{
		"name":    filepath.Base(os.Args[0]),
		"version": version,
		"host":    r.Host,
		"runtime": runtime.Version(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthData)
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

func getFromEnv(value, defaultValue string) string {
	if env, isSet := os.LookupEnv(value); isSet {
		return env
	}
	return defaultValue
}
