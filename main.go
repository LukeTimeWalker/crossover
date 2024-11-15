package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {

	logger := GetLogger()

	filepath := "./config.yml"

	flag.StringVar(&filepath, "config", filepath, "path to the configfile")

	flag.Parse()

	config, err := GetConfig(filepath)

	if err != nil {
		logger.Error(err.Error(), "filepath", filepath)
		os.Exit(1)
	}

	SetLogLevel(config.LogLevel)

	var envs map[string]interface{}

	if config.EnvironementVariables != nil {
		logger.Info("Reading Environment Variables...")
		envs = LoadEnvs(config.EnvironementVariables)
	}

	router := http.NewServeMux()

	for _, route := range config.Routes {
		logger.Info("Registering Route", "name", route.Name, "method", route.Receiver.Method, "route", route.Receiver.Endpoint)

		routeHandler, err := CreateRouteHandler(&route, &envs)

		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		receiverEndpoint := path.Clean(fmt.Sprintf("/%s", route.Receiver.Endpoint))

		router.HandleFunc(fmt.Sprintf("%s %s", route.Receiver.Method, receiverEndpoint), routeHandler)
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
