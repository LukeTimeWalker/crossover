package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	config := GetConfig("./config.yml")

	logger := GetLogger()
	SetLogLevel(config.LogLevel)

	client := http.Client{}

	var envs map[string]interface{}

	if config.EnvironementVariables != nil {
		logger.Info("Reading Environment Variables...")
		envs = LoadEnvs(config.EnvironementVariables)
	}

	router := mux.NewRouter()

	for _, endpoint := range config.Endpoints {
		logger.Info("Registering Endpoint", "method", endpoint.Receiver.Method, "endpoint", endpoint.Receiver.Path)
		RegisterEndpoint(router, &endpoint, &client, &envs)
	}

	logger.Info("Starting Crossover Server...")

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", config.Port),
		WriteTimeout: time.Duration(config.Timeout) * time.Second,
		ReadTimeout:  time.Duration(config.Timeout) * time.Second,
	}

	panic(server.ListenAndServe())
}
