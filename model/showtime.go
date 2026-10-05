package model

import "time"

type Showtime struct {
	ID       uint `gorm:"primaryKey"`
	MovieID  uint
	CinemaID uint
	Showtime time.Time

	Movie  Movie
	Cinema Cinema
}

type MovieShowtimeResponse struct {
	CinemaID  uint                     `json:"cinema_id"`
	Showtimes []CinemaShowtimeResponse `json:"cinema_showtimes"`
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
}
