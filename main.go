package main

import (
	"flag"

	"example/configreader"
	"example/internal/server"
	"example/pkg/config"
	"example/pkg/core"
)

var (
	appName        = "example-service"
	configFileFlag = flag.String("config.file", "", "Path to configuration file.")
)

func main() {
	// parse command-line flags
	flag.Parse()

	config.Init(
		config.New(
			config.WithDefaultEnvVars(""),
			config.WithDefaultConfigFile(appName, *configFileFlag),
		),
	)
	err := config.UnmarshalToStruct("", &configreader.Config)
	if err != nil {
		panic(err)
	}
	a := core.NewApp(
		core.AppWithName(appName),
		core.AppWithLogger(
			configreader.Config.Logger.Mode,
			configreader.Config.Logger.Debug,
			configreader.Config.Logger.Sensitive,
		),
		core.AppWithAction(
			server.NewInstance(),
		),
	)

	a.Run()
}
