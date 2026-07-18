package core

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/LukeTimeWalker/crossover/pkg/utils"
)

type Crossover struct {
	logger *slog.Logger
	router *http.ServeMux
	server *http.Server
	mu     sync.RWMutex
	handlers map[string]RouteHandler
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

	crossover := Crossover{logger: logger, router: router, server: server, handlers: make(map[string]RouteHandler)}

	for _, route := range config.Routes {
		if err := crossover.RegisterRoute(route, envs); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
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

	pattern := routePattern(route)

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.handlers[pattern]; exists {
		c.handlers[pattern] = nil
	}
}

func (c *Crossover) RegisterRoute(route Route, envs map[string]any) error {
	c.logger.Info("Registering Route", "name", route.Name, "method", route.Receiver.Method, "route", route.Receiver.Endpoint)

	routeHandler, err := CreateRouteHandler(&route, &envs)

	if err != nil {
		return err
	}

	pattern := routePattern(route)

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.handlers[pattern]; !exists {
		c.router.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
			c.mu.RLock()
			handler := c.handlers[pattern]
			c.mu.RUnlock()

			if handler == nil {
				http.NotFound(writer, request)
				return
			}

			handler(writer, request)
		})
	}

	c.handlers[pattern] = routeHandler

	return nil
}

func routePattern(route Route) string {
	receiverEndpoint := path.Clean(fmt.Sprintf("/%s", route.Receiver.Endpoint))
	return fmt.Sprintf("%s %s", route.Receiver.Method, receiverEndpoint)
}

func (c *Crossover) Router() *http.ServeMux {
	return c.router
}
