package model

import "time"

type Showtime struct {
	ID       uint `gorm:"primaryKey"`
	MovieID  uint
	CinemaID uint
	Showtime time.Time
}

type ShowtimeResponse struct {
	ID       uint `gorm:"primaryKey"`
	MovieID  uint
	CinemaID uint
	Showtime time.Time
}

type CreateShowtimeRequest struct {
	MovieID  uint      `json:"movie_id"`
	CinemaID uint      `json:"cinema_id"`
	Showtime time.Time `json:"showtime"`
}

type MovieShowtimeResponse struct {
	CinemaID  uint                     `json:"cinema_id"`
	Showtimes []CinemaShowtimeResponse `json:"showtimes"`
}

type ShowtimeSeatsResponse struct {
	CinemaID uint                    `json:"cinema_id"`
	Showtime CinemaShowtimeResponse  `json:"showtime"`
	Seats    []ShowtimeAvailableSeat `json:"seats"`
}

type ShowtimeAvailableSeat struct {
	CinemaSeatResponse
	Available bool `json:"available"`
}

type CinemaShowtimeResponse struct {
	ID       uint      `json:"id"`
	Showtime time.Time `json:"showtime"`
	// StartsAt is the showtime as a real instant in the cinema's time zone;
	// Bookable is false from 15 minutes after it starts.
	StartsAt time.Time `json:"starts_at"`
	Bookable bool      `json:"bookable"`
}
