package server

import (
	"example/internal/handler/personnel"
	"example/internal/storage"
	"example/pkg/core"
	"example/pkg/http"

	"github.com/gin-gonic/gin"
)

func GetRoutes(ds *storage.DataStorage) []core.HTTPRouterOption {
	personnelHandler := personnel.New(ds)
	var handlers = []http.AbstractHandler{
		&http.GroupEndpointHandler{
			GroupPath:       "/personnel",
			GroupMiddleware: []gin.HandlerFunc{SimpleAuthMiddleware},
			EndpointHandlers: []*http.EndpointHandler{
				{
					Endpoint:   "/:personnelId",
					Method:     "GET",
					Middleware: nil,
					Handler:    personnelHandler.GetPersonnel,
				},
				{
					Endpoint:   "",
					Method:     "POST",
					Middleware: []gin.HandlerFunc{AdminMiddleware},
					Handler:    personnelHandler.CreatePersonnel,
				},
				// TODO: adding more function here
				{
					Endpoint:   "",
					Method:     "GET",
					Middleware: []gin.HandlerFunc{},
					Handler:    personnelHandler.ListPersonnel,
				},
				{
					Endpoint:   "/:personnelId",
					Method:     "PUT",
					Middleware: []gin.HandlerFunc{AdminMiddleware},
					Handler:    personnelHandler.UpdatePersonnel,
				},
				{
					Endpoint:   "/:personnelId",
					Method:     "DELETE",
					Middleware: []gin.HandlerFunc{AdminMiddleware},
					Handler:    personnelHandler.DeletePersonnel,
				},
			},
		},
	}

	var options = []core.HTTPRouterOption{}
	for _, h := range handlers {
		options = append(options, h.ToHTTPRouteOptions()...)
	}

	return options
}

func CORS() core.HTTPRouterOption {
	return func(router *core.HttpRouter) {
		router.BaseRoute().Use(func(c *gin.Context) {
			origin := c.GetHeader("origin")
			if origin == "" {
				origin = "*"
			}
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})
	}
}
func Recover() core.HTTPRouterOption {
	return func(router *core.HttpRouter) {
		router.BaseRoute().Use(gin.Recovery())
	}
}
