package main

import (
	"flag"
	"log"
	"time"

	"github.com/LukeTimeWalker/crossover/pkg/core"
)

func main() {

	filepath := "./config.yaml"

	flag.StringVar(&filepath, "config", filepath, "path to the config file")

	flag.Parse()

	config, err := core.GetConfig(filepath)

	if err != nil {
		log.Fatal(err.Error(), "filepath", filepath)
	}

	crossover := core.NewCrossover(config)

	go crossover.Start()

	time.Sleep(10 * time.Second)

	crossover.RegisterRoute(core.Route{Name: "dynamic", Receiver: core.Receiver{Endpoint: "dynamic", Method: "GET"}}, nil)

	time.Sleep(5 * time.Second)

	crossover.UnregisterRoute(core.Route{Name: "dynamic", Receiver: core.Receiver{Endpoint: "dynamic", Method: "GET"}})

	select {}

}
