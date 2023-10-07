package main

import (
	"flag"

	"github.com/klyngen/votomatic-3000/packages/backend/presentation"
	"github.com/klyngen/votomatic-3000/packages/backend/utils"
)

var configFileName string
var websocketPort int
var restPort int

func readCommandLineArguments() {
	flag.IntVar(&restPort, "restport", 8080, "Specify a port")
	flag.StringVar(&configFileName, "config", "./config.json", "Specify a file with questions")

	flag.Parse()
}

func main() {
	readCommandLineArguments()

	if len(configFileName) == 0 {
		flag.Usage()
		return
	}

	configuration, err := utils.ReadConfiguration(configFileName)

	if err != nil {
		panic(err.Error())
	}

	if !configuration.IsValidConfiguration() {
		panic("Not a valid configuration")
	}

	api := presentation.NewPresentationAPI(*configuration)

	api.ListenAndServe(restPort)

}
