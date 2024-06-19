package server

import "github.com/gin-gonic/gin"

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
