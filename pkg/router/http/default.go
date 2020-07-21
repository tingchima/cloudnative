package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisteDefault ...
func RegisteDefault(g *gin.Engine) *gin.Engine {
	g.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "["+c.ClientIP()+"] pong!")
		return
	})
	g.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome gin server...",
		})
	})
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "unmatched router",
		})
	})

	return g
}
