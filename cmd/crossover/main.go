package main

import (
	"flag"
	"log"

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



	select {}

}
