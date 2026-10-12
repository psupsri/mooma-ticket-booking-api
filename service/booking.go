package service

import (
	"encoding/json"
	"errors"
	"mooma-api/model"
	"mooma-api/repository"
	"mooma-api/utils"
)

var (
	ErrNoSeats            = errors.New("at least one seat is required")
	ErrSeatAlreadyBooked  = repository.ErrSeatAlreadyBooked
	ErrBookingNotFound    = repository.ErrBookingNotFound
	ErrCancellationClosed = errors.New("bookings can only be cancelled up to 2 hours before the showtime")
	ErrShowtimeNotFound   = repository.ErrShowtimeNotFound
	ErrBookingClosed      = errors.New("this showtime is no longer open for booking")
)

type BookingService struct {
	repository *repository.BookingRepository
	schedule   *Schedule
}

func NewBookingService(repository *repository.BookingRepository, schedule *Schedule) *BookingService {
	return &BookingService{
		repository: repository,
		schedule:   schedule,
	}
}

func (s *BookingService) CreateBooking(request model.CreateBookingRequest, customerID uint) (model.GetBookingResponse, error) {
	var response model.GetBookingResponse
	if len(request.Seats) == 0 {
		return response, ErrNoSeats
	}
	seatsJSON, err := json.Marshal(request.Seats)
	if err != nil {
		return response, err
	}
	booking := model.Booking{
		CustomerID: customerID,
		ShowtimeID: request.ShowtimeID,
		Seats:      seatsJSON,
	}

	err = s.repository.CreateBooking(&booking, request.Seats, func(showtime model.Showtime) error {
		if !s.schedule.BookingOpen(showtime) {
			return ErrBookingClosed
		}
		return nil
	})
	if err != nil {
		return response, err
	}

	created, movie, err := s.repository.GetByIDWithDetails(booking.ID)
	if err != nil {
		return response, err
	}
	return s.toBookingResponse(created, movie)
}

func (s *BookingService) GetCustomerBookings(customerID uint) ([]model.GetBookingResponse, error) {
	var response = []model.GetBookingResponse{}
	bookings, movies, err := s.repository.GetByCustomerID(customerID)
	if err != nil {
		return response, err
	}
	for _, booking := range bookings {
		item, err := s.toBookingResponse(booking, movies[booking.Showtime.MovieID])
		if err != nil {
			return response, err
		}
		response = append(response, item)
	}
	return response, nil
}

func (s *BookingService) CancelBooking(id string, customerID uint) error {
	bookingID, err := utils.StringToUint(id)
	if err != nil {
		return ErrBookingNotFound
	}
	return s.repository.Delete(bookingID, customerID, func(booking model.Booking) error {
		if !s.schedule.CanCancel(booking.Showtime) {
			return ErrCancellationClosed
		}
		return nil
	})
}

func (s *BookingService) toBookingResponse(booking model.Booking, movie model.Movie) (model.GetBookingResponse, error) {
	var response model.GetBookingResponse
	seats, err := toCinemaSeatsResponse(booking.Seats)
	if err != nil {
		return response, err
	}

	return model.GetBookingResponse{
		ID: booking.ID,
		Customer: model.CustomerResponse{
			FirstName: booking.Customer.FirstName,
			LastName:  booking.Customer.LastName,
			Tel:       booking.Customer.Tel,
			Email:     booking.Customer.Email,
		},
		Movie: model.MovieResponse{
			ID:        movie.ID,
			Title:     movie.Title,
			PosterURL: movie.PosterURL,
		},
		CinemaID:       booking.Showtime.CinemaID,
		Showtime:       s.schedule.showtimeResponse(booking.Showtime),
		Seats:          seats,
		Status:         booking.Status,
		CreatedAt:      booking.CreatedAt,
		StartsAt:       s.schedule.StartsAt(booking.Showtime),
		CancelDeadline: s.schedule.CancelDeadline(booking.Showtime),
		CanCancel:      s.schedule.CanCancel(booking.Showtime),
	}, nil
}
