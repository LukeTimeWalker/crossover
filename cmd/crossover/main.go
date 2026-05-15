package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/LukeTimeWalker/crossover/pkg/core"
	"github.com/LukeTimeWalker/crossover/pkg/utils"
)

func main() {

	logger := utils.GetLogger()

	filepath := "./config.yaml"

	flag.StringVar(&filepath, "config", filepath, "path to the config file")

	flag.Parse()

	config, err := core.GetConfig(filepath)

	if err != nil {
		logger.Error(err.Error(), "filepath", filepath)
		os.Exit(1)
	}

	utils.SetLogLevel(config.LogLevel)

	var envs map[string]any

	if config.EnvironementVariables != nil {
		logger.Info("Reading Environment Variables...")
		envs = utils.LoadEnvs(config.EnvironementVariables)
	}

	router := http.NewServeMux()

	core.InitializeClients(getNumberOfSenders(config.Routes))

	for _, route := range config.Routes {
		logger.Info("Registering Route", "name", route.Name, "method", route.Receiver.Method, "route", route.Receiver.Endpoint)

		routeHandler, err := core.CreateRouteHandler(&route, &envs)

		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		receiverEndpoint := path.Clean(fmt.Sprintf("/%s", route.Receiver.Endpoint))

		router.HandleFunc(fmt.Sprintf("%s %s", route.Receiver.Method, receiverEndpoint), routeHandler)
	}

	logger.Info("Starting Crossover Server...")

	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", config.Port),
	}

	panic(server.ListenAndServe())
}

func getNumberOfSenders(routes []core.Route) int {
	numberOfSender := 0
	for i := range routes {
		numberOfSender += len(routes[i].Sender)
	}
	return numberOfSender
}
