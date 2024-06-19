package core

import (
	"context"
	"time"

	"example/pkg/config"
	"example/pkg/log"
)

type AppOption func(a *App)

// AppWithName set app's name
func AppWithName(name string) AppOption {
	return func(a *App) {
		a.name = name
	}
}

// AppWithContext set contect for app
func AppWithContext(ctx context.Context) AppOption {
	return func(a *App) {
		a.baseCtx = ctx
	}
}

// AppWithLogger init log package instance
func AppWithLogger(mode string, debug bool, sensitiveData bool) AppOption {
	return func(a *App) {
		log.Init(
			log.New(
				log.WithModeFromString(mode),
				log.WithDebug(debug),
				log.EnableSensitive(sensitiveData),
			),
		)
	}
}

// Deprecated: should init config as first step before create app with AppOption
//
// AppWithConfig init config from environment prefix envPrefix and configFile
func AppWithConfig(envPrefix, configFile string) AppOption {
	return func(a *App) {
		config.Init(
			config.New(
				config.WithDefaultEnvVars(envPrefix),
				config.WithDefaultConfigFile(a.name, configFile),
			),
		)
	}
}

// AppWithAction set execute action (start/stop) for process's lifecycle
func AppWithAction(actions AppAction) AppOption {
	return func(a *App) {
		a.startFunc = actions.Start
		a.shutdownFunc = actions.Shutdown
	}
}

// AppWithGracefulTimeout setup graceful timeout
func AppWithGracefulTimeout(waitTime time.Duration) AppOption {
	return func(a *App) {
		a.gracefulTimeout = waitTime
	}
}
