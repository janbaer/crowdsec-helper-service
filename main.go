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
	mux.HandleFunc("DELETE /crowdsec-helper-service/decisions", handleDeleteCrowdsecDecison)
	mux.HandleFunc("POST /crowdsec-helper-service/decisions", handlePostCrowdsecDecison)

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
	ipAddress := queryParams.Get("ip")
	if net.ParseIP(ipAddress) == nil {
		logger.Error("The provided IP address is not valid", "ipAddress", ipAddress)
		sendResponse(w, http.StatusBadRequest, "The provided IP address is not valid\n")
		return
	}

	if err := csclirunner.DeleteDecision(ipAddress); err != nil {
		sendResponse(w, http.StatusInternalServerError, "Deletion of the decision failed\n")
		return
	}

	sendResponse(w, http.StatusOK, "OK")
}

func handlePostCrowdsecDecison(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	ipAddress := queryParams.Get("ip")
	decisonType := queryParams.Get("type")
	duration := queryParams.Get("duration")

	if net.ParseIP(ipAddress) == nil {
		logger.Error("The provided IP address is not valid", "ipAddress", ipAddress)
		sendResponse(w, http.StatusBadRequest, "The provided IP address is not valid\n")
		return
	}

	if decisonType != "ban" && decisonType != "captcha" {
		logger.Error("The provided decision type is not valid", "type", decisonType)
		sendResponse(w, http.StatusBadRequest, "The provided decision type is not valid\n")
		return
	}

	if len(duration) == 0 {
		logger.Error("The provided duration is not valid", "duration", duration)
		sendResponse(w, http.StatusBadRequest, "The provided duration is not valid\n")
		return
	}

	if err := csclirunner.CreateDecision(ipAddress, decisonType, duration); err != nil {
		sendResponse(w, http.StatusInternalServerError, "Creation of a new decision failed\n")
		return
	}

	sendResponse(w, http.StatusCreated, "Created")
}

func getFromEnv(value, defaultValue string) string {
	if env, isSet := os.LookupEnv(value); isSet {
		return env
	}
	return defaultValue
}
