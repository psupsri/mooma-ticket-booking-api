package model

import (
	"time"

	"gorm.io/datatypes"
)

type BookingStatus string

const (
	Confirmed BookingStatus = "confirmed"
	Cancelled BookingStatus = "cancelled"
	Expired   BookingStatus = "expired"
)

type Booking struct {
	ID         uint `gorm:"primaryKey"`
	CustomerID uint
	ShowtimeID uint
	Seats      datatypes.JSON
	Status     BookingStatus `gorm:"not null;default:confirmed"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Customer Customer
	Showtime Showtime
}

type GetBookingResponse struct {
	ID       uint                   `json:"id"`
	Customer CustomerResponse       `json:"customer"`
	Movie    MovieResponse          `json:"movie"`
	CinemaID uint                   `json:"cinema_id"`
	Showtime CinemaShowtimeResponse `json:"showtime"`
	Seats    []CinemaSeatResponse   `json:"seats"`
}

type BookedShowtimeSeats struct {
	ShowtimeID uint   `json:"showtime_id"`
	BookingID  uint   `json:"booking_id"`
	SeatLine   string `json:"seat_line"`
	SeatNumber string `json:"seat_number"`
}
