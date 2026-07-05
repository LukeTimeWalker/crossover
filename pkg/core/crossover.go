package core

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"

	"github.com/LukeTimeWalker/crossover/pkg/utils"
)

type Crossover struct {
	logger *slog.Logger
	router *http.ServeMux
	server *http.Server
}

func NewCrossover(config *Config) *Crossover {
	logger := utils.GetLogger()
	utils.SetLogLevel(config.LogLevel)

	var envs map[string]any

	if config.EnvironementVariables != nil {
		logger.Info("Reading Environment Variables...")
		envs = utils.LoadEnvs(config.EnvironementVariables)
	}

	router := http.NewServeMux()

	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", config.Port),
	}

	crossover := Crossover{logger: logger, router: router, server: server}

	for _, route := range config.Routes {
		crossover.RegisterRoute(route, envs)
	}

	return &crossover
}

func (c *Crossover) Start() {
	c.logger.Info("Starting Crossover")
	panic(c.server.ListenAndServe())
}

func (c *Crossover) Stop() {
	c.logger.Info("Gracefully shutting down Crossover")
	c.server.Close()
}

func (c *Crossover) UnregisterRoute(route Route) {
	c.logger.Info("Unregistering Route", "name", route.Name, "method", route.Receiver.Method, "route", route.Receiver.Endpoint)

	receiverEndpoint := path.Clean(fmt.Sprintf("/%s", route.Receiver.Endpoint))

	c.router.HandleFunc(fmt.Sprintf("%s %s", route.Receiver.Method, receiverEndpoint), func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}

func (c *Crossover) RegisterRoute(route Route, envs map[string]any) {
	c.logger.Info("Registering Route", "name", route.Name, "method", route.Receiver.Method, "route", route.Receiver.Endpoint)

	routeHandler, err := CreateRouteHandler(&route, &envs)

	if err != nil {
		c.logger.Error(err.Error())
		os.Exit(1)
	}

	receiverEndpoint := path.Clean(fmt.Sprintf("/%s", route.Receiver.Endpoint))

	c.router.HandleFunc(fmt.Sprintf("%s %s", route.Receiver.Method, receiverEndpoint), routeHandler)
}

func (c *Crossover) Router() *http.ServeMux {
	return c.router
}
