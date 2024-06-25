package main

import (
	"flag"

	"example/configreader"
	"example/internal/mqserver"
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
	var appAction core.AppAction
	if cfg.RunMode == "api" {
		appAction = server.NewInstance()
	} else {
		appAction = mqserver.NewInstance()
	}
	a := core.NewApp(
		core.AppWithName(appName),
		core.AppWithLogger(
			cfg.Logger.Mode,
			cfg.Logger.Debug,
			cfg.Logger.Sensitive,
		),
		core.AppWithAction(
			appAction,
		),
	)

	a.Run()
}
