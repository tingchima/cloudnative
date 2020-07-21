package http

import "github.com/gin-gonic/gin"

// RegisteAuth ...
func RegisteAuth(g *gin.Engine) *gin.Engine {

	g.POST("/auth")

	return g
}
