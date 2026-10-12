package config

import "os"

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	JWTSecret        string
	// IANA zone of the cinemas. Showtimes are stored as local wall-clock
	// times (TIMESTAMP without time zone) and are read in this zone.
	CinemaTimezone string
}

const defaultCinemaTimezone = "Asia/Bangkok"

func Load() Config {
	return Config{
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		CinemaTimezone:   getEnv("CINEMA_TIMEZONE", defaultCinemaTimezone),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
