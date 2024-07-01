package http

import (
	"example/pkg/core"

	"github.com/gin-gonic/gin"
)

type ServiceHealth interface {
	Name() string
	Check() error
}

func getHealthCheckFunc(services ...ServiceHealth) gin.HandlerFunc {
	return func(context *gin.Context) {
		status := make(map[string]string)
		isHealth := true
		for _, s := range services {
			err := s.Check()
			if err != nil {
				status[s.Name()] = err.Error()
				isHealth = false
			} else {
				status[s.Name()] = "healthy"
			}
		}
		x := map[string]interface{}{
			"isHealth": isHealth,
			"status":   status,
		}
		httpCode := 200
		if !isHealth {
			httpCode = 500
		}
		context.JSON(httpCode, x)
	}
}

func RegisterHealthCheck(path string, services ...ServiceHealth) core.HTTPRouterOption {
	return func(router *core.HttpRouter) {
		router.BaseRoute().GET(path, getHealthCheckFunc(services...))
	}
}
