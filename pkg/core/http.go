package core

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"example/pkg/log"

	"github.com/gin-gonic/gin"
)

var (
	ErrServerClosed = http.ErrServerClosed
)

type HTTPServer struct {
	name              string
	address           string
	server            *http.Server
	handler           http.Handler
	tlsConfig         *tls.Config
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
}

func NewHTTPServer(options ...HTTPServerOption) *HTTPServer {
	server := &HTTPServer{}

	for _, option := range options {
		option(server)
	}

	if server.name == "" {
		log.Panic("must set name for HTTP server")
	}

	return server
}

// Run serve HTTP server
func (s *HTTPServer) Run() error {

	s.server = &http.Server{
		Addr:              s.address,
		Handler:           s.handler,
		TLSConfig:         s.tlsConfig,
		ReadTimeout:       s.readTimeout,
		ReadHeaderTimeout: s.readHeaderTimeout,
		WriteTimeout:      s.writeTimeout,
		IdleTimeout:       s.idleTimeout,
	}
	log.Infof("HTTP %v listening on: %v", s.name, s.address)
	return s.server.ListenAndServe()
}

// Stop gracefully shutdown HTTP server
func (s *HTTPServer) Stop(ctx context.Context) {
	log.Infof("HTTP %v shutting down...", s.name)
	if s.server != nil {
		err := s.server.Shutdown(ctx)
		if err != nil {
			log.Error(err)
		}
	}
	log.Infof("HTTP %v gracefully stopped", s.name)
}

type HttpRouter struct {
	base *gin.Engine
}

func (router *HttpRouter) BaseRoute() *gin.Engine {
	return router.base
}

type HTTPRouterOption func(router *HttpRouter)

func NewHTTPRouter(opts ...HTTPRouterOption) *gin.Engine {
	router := &HttpRouter{
		base: gin.New(),
	}

	for _, opt := range opts {
		opt(router)
	}

	return router.base
}

func HTTPRouterWithNoFoundHandler(handler gin.HandlerFunc) HTTPRouterOption {
	return func(router *HttpRouter) {
		router.base.NoRoute(handler)
	}
}

func HTTPRouterWithMethodNotAllowedHandler(handler gin.HandlerFunc) HTTPRouterOption {
	return func(router *HttpRouter) {
		router.base.NoMethod(handler)
	}
}
