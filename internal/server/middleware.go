package server

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func SimpleAuthMiddleware(c *gin.Context) {
	if c.GetHeader("Authorization") != "Bearer 123" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

func AdminMiddleware(c *gin.Context) {
	if c.GetHeader("X-ADMIN") != "true" {
		c.JSON(403, gin.H{"error": "Forbidden"})
		c.Abort()
		return
	}
	c.Next()
}

// Setups the Prometheus cache metrics, can change if needed
var prometheusCacheCounting = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "cache_metrics",
		Help: "Counting cache hit and miss",
	}, []string{"state"})

func CacheMissFunc() {
	prometheusCacheCounting.WithLabelValues("cache_miss").Inc()
}

func CacheHitFunc() {
	prometheusCacheCounting.WithLabelValues("cache_hit").Inc()
}

func TotalCountFunc() {
	prometheusCacheCounting.WithLabelValues("total").Inc()
}
