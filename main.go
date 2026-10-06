package main

import (
	"log"
	"mooma-api/config"
	"mooma-api/handlers"
	"mooma-api/repository"
	"mooma-api/routes"
	"mooma-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	c := config.Load()
	db, err := config.ConnectDatabase(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	appHandlers := routes.Handlers{
		AuthHandler: authHandler,
	}

	routes.RegisterRoutes(r, appHandlers)

	r.Run()
}
