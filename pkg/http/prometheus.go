package http

import (
	"context"

	"example/pkg/core"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	prometheusmiddleware "github.com/slok/go-http-metrics/middleware"
)

type reporter struct {
	c *gin.Context
}

func (r *reporter) Method() string { return r.c.Request.Method }

func (r *reporter) Context() context.Context { return r.c.Request.Context() }

func (r *reporter) URLPath() string { return r.c.FullPath() }

func (r *reporter) StatusCode() int { return r.c.Writer.Status() }

func (r *reporter) BytesWritten() int64 { return int64(r.c.Writer.Size()) }

func Prometheus() core.HTTPRouterOption {
	return func(router *core.HttpRouter) {
		prometheusMiddleware := prometheusmiddleware.New(prometheusmiddleware.Config{
			Recorder: metrics.NewRecorder(metrics.Config{}),
		})
		router.BaseRoute().Use(func(c *gin.Context) {
			r := &reporter{c: c}
			prometheusMiddleware.Measure("", r, func() {
				c.Next()
			})
		})
	}
}

func StartPrometheus(host string, port int) *core.HTTPServer {
	prometheusServer := core.NewHTTPServer(
		core.HTTPServerWithName("prometheus"),
		core.HTTPServerWithAddress(
			host,
			port,
		),
		core.HTTPServerWithHandler(promhttp.Handler()))
	go prometheusServer.Run()

	return prometheusServer
}
