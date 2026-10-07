package main

import (
	"log"
	"mooma-api/config"
	"mooma-api/handlers"
	"mooma-api/repository"
	"mooma-api/routes"
	"mooma-api/service"

	"github.com/gin-gonic/gin"
)

func main() {
	c := config.Load()
	db, err := config.ConnectDatabase(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := gin.Default()

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	movieRepo := repository.NewMovieRepository(db)
	movieService := service.NewMovieService(movieRepo)
	movieHandler := handlers.NewMovieHandler(movieService)

	showtimeRepo := repository.NewShowtimeRepository(db)
	showtimeService := service.NewShowtimeService(showtimeRepo)
	showtimeHandler := handlers.NewShowtimeHandler(showtimeService)

	appHandlers := routes.Handlers{
		AuthHandler:     authHandler,
		MovieHandler:    movieHandler,
		ShowtimeHandler: showtimeHandler,
	}

	routes.RegisterRoutes(r, appHandlers)

	r.Run()
}
