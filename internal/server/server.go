package server

import (
	"context"

	"example/configreader"
	"example/internal/storage"
	"example/pkg/cache"
	"example/pkg/cache/redis_cache"
	"example/pkg/core"
	"example/pkg/driver/postgresql"
	"example/pkg/driver/redis"
	"example/pkg/http"
	"example/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// Instance represents an instance of the server
type Instance struct {
	ctx              context.Context
	httpServer       *core.HTTPServer
	prometheusServer *core.HTTPServer
	ds               *storage.DataStorage
	cache            cache.Cache
}

// NewInstance returns a new instance of our server
func NewInstance() *Instance {
	return &Instance{}
}

// Start starts the server
func (i *Instance) Start(ctx context.Context) {
	i.ctx = ctx
	if log.Env() != log.ModeDevelopment {
		gin.SetMode(gin.ReleaseMode)
	}
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
	}, configreader.Config.Postgresql.LogSQL)
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
	// init redis standalone, see more type of redis in pkg/driver/redis
	redisClient, err := redis.NewConnection(&redis.SingleConnection{
		Address:    configreader.Config.Redis.Address,
		Password:   configreader.Config.Redis.Password,
		DB:         configreader.Config.Redis.DB,
		MaxRetries: configreader.Config.Redis.MaxRetries,
		PoolSize:   configreader.Config.Redis.PoolSize,
	})
	if err != nil {
		log.Fatalf("init redis error", "error", err)
	}
	redisCache, err := redis_cache.New(redis_cache.WithName("example-cache"),
		redis_cache.WithContext(ctx),
		redis_cache.WithClient(redisClient))
	if err != nil {
		log.Warnw("init cache error", "error", err)
	}
	if configreader.Config.Prometheus.Enable {
		i.cache = cache.NewMetricAbleCache(redisCache, CacheHitFunc, CacheMissFunc, TotalCountFunc)
	} else {
		i.cache = redisCache
	}

	if i.cache != nil {
		// set cache for data storage
		i.ds.SetCache(i.cache)
	}
	var routerOption []core.HTTPRouterOption
	routerOption = append(routerOption, Recover())
	routerOption = append(routerOption, CORS())
	routerOption = append(routerOption, http.RegisterHealthCheck("/healthz", i.ds))
	if configreader.Config.Prometheus.Enable {
		routerOption = append(routerOption, http.Prometheus())
		i.prometheusServer = http.StartPrometheus(configreader.Config.Prometheus.Host, configreader.Config.Prometheus.Port)
	}

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
	if i.prometheusServer != nil {
		i.prometheusServer.Stop(ctx)
	}
}
