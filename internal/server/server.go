package server

import (
	"context"

	"example/configreader"
	"example/internal/storage"
	"example/pkg/core"
	"example/pkg/driver/postgresql"
	"example/pkg/log"

	"github.com/spf13/cast"
)

// Instance represents an instance of the server
type Instance struct {
	ctx        context.Context
	httpServer *core.HTTPServer
	ds         *storage.DataStorage
}

// NewInstance returns a new instance of our server
func NewInstance() *Instance {
	return &Instance{}
}

// Start starts the server
func (i *Instance) Start(ctx context.Context) {
	i.ctx = ctx
	var err error
	db, err := postgresql.New(postgresql.Connection{
		Host:         configreader.Config.Postgresql.Host,
		Port:         configreader.Config.Postgresql.Port,
		Username:     configreader.Config.Postgresql.Username,
		Password:     configreader.Config.Postgresql.Password,
		DatabaseName: configreader.Config.Postgresql.DatabaseName,
		Schema:       configreader.Config.Postgresql.Schema,
		MaxIdleConns: configreader.Config.Postgresql.MaxIdleConns,
		MaxOpenConns: configreader.Config.Postgresql.MaxOpenConns,
		MaxLifetime:  cast.ToDuration(configreader.Config.Postgresql.MaxLifetime),
	})
	if err != nil {
		log.Fatalf("init db error", "error", err)
	}
	i.ds = storage.New(db)
	if configreader.Config.Postgresql.AutoMigrate {
		// TODO: don't using migrate in production.
		err := i.ds.Migrate()
		if err != nil {
			log.Fatalf("migrate db error", "error", err)
		}
	}
	var routerOption []core.HTTPRouterOption
	routerOption = append(routerOption, Recover())
	routerOption = append(routerOption, CORS())

	routerOption = append(routerOption, GetRoutes(i.ds)...)
	router := core.NewHTTPRouter(
		routerOption...,
	)

	// Startup the HTTP HTTPServer in a way that we can gracefully shut it down again
	i.httpServer = core.NewHTTPServer(
		core.HTTPServerWithName("chain-admin-server"),
		core.HTTPServerWithAddress(
			configreader.Config.HTTP.Host,
			configreader.Config.HTTP.Port,
		),
		core.HTTPServerWithHandler(router))
	err = i.httpServer.Run()
	if err != core.ErrServerClosed {
		log.Errorw("HTTP HTTPServer stopped unexpected", "error", err)
		i.Shutdown(ctx)
	}
}

// Shutdown stops the server
func (i *Instance) Shutdown(ctx context.Context) {
	// Shutdown HTTP server
	i.httpServer.Stop(ctx)
}
