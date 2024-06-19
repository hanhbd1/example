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
	cfg := configreader.Config
	if err != nil {
		panic(err)
	}

	a := core.NewApp(
		core.AppWithName(appName),
		core.AppWithLogger(
			cfg.Logger.Mode,
			cfg.Logger.Debug,
			cfg.Logger.Sensitive,
		),
		core.AppWithAction(
			server.NewInstance(),
		),
	)

	a.Run()
}
