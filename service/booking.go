package service

import (
	"encoding/json"
	"mooma-api/model"
	"mooma-api/repository"
)

type BookingService struct {
	repository *repository.BookingRepository
}

func NewBookingService(repository *repository.BookingRepository) *BookingService {
	return &BookingService{
		repository: repository,
	}
}

func (s *BookingService) CreateBooking(request model.CreateBookingRequest, customerID uint) (model.GetBookingResponse, error) {
	var response model.GetBookingResponse
	seatsJSON, err := json.Marshal(request.Seats)
	if err != nil {
		return response, err
	}
	booking := model.Booking{
		CustomerID: customerID,
		ShowtimeID: request.ShowtimeID,
		Seats:      seatsJSON,
	}

	if err := s.repository.CreateBooking(&booking, request.Seats); err != nil {
		return response, nil
	}

	return response, nil
}
