package http

import (
	"fmt"

	"example/pkg/core"

	"github.com/gin-gonic/gin"
)

type AbstractHandler interface {
	ToHTTPRouteOptions() []core.HTTPRouterOption
}

type EndpointHandler struct {
	Endpoint   string
	Method     string
	Middleware []gin.HandlerFunc
	Handler    gin.HandlerFunc
}

func (e *EndpointHandler) ToHTTPRouteOptions() []core.HTTPRouterOption {
	x := []core.HTTPRouterOption{
		func(router *core.HttpRouter) {
			var h []gin.HandlerFunc
			if len(e.Middleware) > 0 {
				h = append(h, e.Middleware...)
			}
			h = append(h, e.Handler)
			router.BaseRoute().Handle(e.Method, e.Endpoint, h...)
		},
	}
	return x
}

type GroupEndpointHandler struct {
	GroupPath        string
	GroupMiddleware  []gin.HandlerFunc
	EndpointHandlers []*EndpointHandler
}

func (e *GroupEndpointHandler) ToHTTPRouteOptions() []core.HTTPRouterOption {
	x := []core.HTTPRouterOption{}
	for _, endpoint := range e.EndpointHandlers {
		var h []gin.HandlerFunc
		if len(e.GroupMiddleware) > 0 {
			h = append(h, e.GroupMiddleware...)
		}
		if len(endpoint.Middleware) > 0 {
			h = append(h, endpoint.Middleware...)
		}
		endpoint.Middleware = h
		endpoint.Endpoint = fmt.Sprintf("%s%s", e.GroupPath, endpoint.Endpoint)
		x = append(x, endpoint.ToHTTPRouteOptions()...)
	}
	return x
}
