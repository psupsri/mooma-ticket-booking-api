package main

import (
	"log"
	"mooma-api/config"
	"mooma-api/handlers"
	"mooma-api/repository"
	"mooma-api/routes"
	"mooma-api/service"
	"time"
	_ "time/tzdata"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	c := config.Load()
	// An empty secret lets anyone sign valid tokens for any customer.
	if c.JWTSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	cinemaLocation, err := time.LoadLocation(c.CinemaTimezone)
	if err != nil {
		log.Fatalf("invalid CINEMA_TIMEZONE %q: %v", c.CinemaTimezone, err)
	}
	schedule := service.NewSchedule(cinemaLocation)
	db, err := config.ConnectDatabase(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := gin.Default()
	registerApp(r, db, schedule)
	r.Run(":3000")
}

// registerApp wires repositories, services and handlers onto r. Split from
// main so tests can build the same router against a test database.
func registerApp(r *gin.Engine, db *gorm.DB, schedule *service.Schedule) {
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	showtimeRepo := repository.NewShowtimeRepository(db)
	showtimeService := service.NewShowtimeService(showtimeRepo)
	showtimeHandler := handlers.NewShowtimeHandler(showtimeService)

	movieRepo := repository.NewMovieRepository(db)
	movieService := service.NewMovieService(movieRepo, showtimeRepo, schedule)
	movieHandler := handlers.NewMovieHandler(movieService)

	cinemaRepo := repository.NewCinemaRepository(db)
	cinemaService := service.NewCinemaService(cinemaRepo)
	cinemaHandler := handlers.NewCinemaHandler(cinemaService)

	bookingRepo := repository.NewBookingRepository(db)
	bookingService := service.NewBookingService(bookingRepo, schedule)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	appHandlers := routes.Handlers{
		AuthHandler:     authHandler,
		MovieHandler:    movieHandler,
		ShowtimeHandler: showtimeHandler,
		CinemaHandler:   cinemaHandler,
		BookingHandler:  bookingHandler,
	}

	routes.RegisterRoutes(r, appHandlers)
}
