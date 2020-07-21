package http

import (
	v1 "apigateway/pkg/router/http/v1"
	"apigateway/pkg/service"

	"github.com/gin-gonic/gin"
)

// Handler Restful api handler
type Handler struct {
	Svc service.IService
}

// NewHandler Create restful api handler
func NewHandler(svc service.IService) *Handler {
	return &Handler{
		Svc: svc,
	}
}

// RegisteRouter ...
func RegisteRouter(router *gin.Engine) {
	Scopes(
		router,
		RegisteDefault,
		RegisteAuth,
		v1.RegisteBook,
		// add new http router at here
	)
}

// Scopes To register any route
func Scopes(r *gin.Engine, funcs ...func(*gin.Engine) *gin.Engine) *gin.Engine {
	for _, f := range funcs {
		r = f(r)
	}

	return r
}
