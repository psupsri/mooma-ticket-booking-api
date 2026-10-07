package routes

import (
	"mooma-api/handlers"
	"mooma-api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	AuthHandler     *handlers.AuthHandler
	MovieHandler    *handlers.MovieHandler
	ShowtimeHandler *handlers.ShowtimeHandler
	CinemaHandler   *handlers.CinemaHandler
}

func RegisterRoutes(r *gin.Engine, handler Handlers) {
	AuthRoutes(r, handler.AuthHandler)
	MovieRoutes(r, handler.MovieHandler)
	ShowtimeRoutes(r, handler.ShowtimeHandler)
	CinemaRoutes(r, handler.CinemaHandler)

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
	movieRouter.POST("/", middleware.AuthMiddleware, handler.CreateMovie)
	movieRouter.GET("/", handler.GetAllMovie)
	movieRouter.GET("/:id", handler.GetByID)
}

func ShowtimeRoutes(r *gin.Engine, handler *handlers.ShowtimeHandler) {
	showtimeRouter := r.Group("/showtime")
	showtimeRouter.POST("/", middleware.AuthMiddleware, handler.CreateShowtime)
	showtimeRouter.GET("/", handler.GetAllShowtime)
}

func CinemaRoutes(r *gin.Engine, handler *handlers.CinemaHandler) {
	cinemaRouter := r.Group("/cinema")
	cinemaRouter.GET("/:id", handler.GetByID)
}
