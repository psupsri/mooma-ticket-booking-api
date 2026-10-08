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

	showtimeRepo := repository.NewShowtimeRepository(db)
	showtimeService := service.NewShowtimeService(showtimeRepo)
	showtimeHandler := handlers.NewShowtimeHandler(showtimeService)

	movieRepo := repository.NewMovieRepository(db)
	movieService := service.NewMovieService(movieRepo, showtimeRepo)
	movieHandler := handlers.NewMovieHandler(movieService)

	cinemaRepo := repository.NewCinemaRepository(db)
	cinemaService := service.NewCinemaService(cinemaRepo)
	cinemaHandler := handlers.NewCinemaHandler(cinemaService)

	appHandlers := routes.Handlers{
		AuthHandler:     authHandler,
		MovieHandler:    movieHandler,
		ShowtimeHandler: showtimeHandler,
		CinemaHandler:   cinemaHandler,
	}

	routes.RegisterRoutes(r, appHandlers)

	r.Run(":3000")
}
