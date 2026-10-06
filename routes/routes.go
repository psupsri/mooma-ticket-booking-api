package routes

import (
	"mooma-api/handlers"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	AuthHandler *handlers.AuthHandler
}

func RegisterRoutes(r *gin.Engine, handler Handlers) {
	AuthRoutes(r, handler.AuthHandler)
}

func AuthRoutes(r *gin.Engine, handler *handlers.AuthHandler) {
	authRouter := r.Group("/auth")
	authRouter.POST("/register", handler.Register)
}
