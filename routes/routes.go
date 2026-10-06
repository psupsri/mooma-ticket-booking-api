package routes

import (
	"mooma-api/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	AuthHandler  *handlers.AuthHandler
	MovieHandler *handlers.MovieHandler
}

func RegisterRoutes(r *gin.Engine, handler Handlers) {
	AuthRoutes(r, handler.AuthHandler)
	MovieRoutes(r, handler.MovieHandler)

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
}

func AuthRoutes(r *gin.Engine, handler *handlers.AuthHandler) {
	authRouter := r.Group("/auth")
	authRouter.POST("/register", handler.Register)
	authRouter.POST("/login", handler.Login)
}

func MovieRoutes(r *gin.Engine, handler *handlers.MovieHandler) {
	movieRouter := r.Group("/movie")
	movieRouter.GET("/", handler.GetAllMovie)
}
