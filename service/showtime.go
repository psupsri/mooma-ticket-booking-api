package service

import (
	"mooma-api/model"
	"mooma-api/repository"
)

type ShowtimeService struct {
	repository *repository.ShowtimeRepository
}

func NewShowtimeService(repository *repository.ShowtimeRepository) *ShowtimeService {
	return &ShowtimeService{
		repository: repository,
	}
}

func (s *ShowtimeService) CreateShowtime(request model.CreateShowtimeRequest) (model.ShowtimeResponse, error) {
	var response model.ShowtimeResponse
	showtime := model.Showtime{
		MovieID:  request.MovieID,
		CinemaID: request.CinemaID,
		Showtime: request.Showtime,
	}
	if err := s.repository.CreateShowtime(&showtime); err != nil {
		return response, err
	}
	response = model.ShowtimeResponse{
		ID:       showtime.ID,
		MovieID:  showtime.MovieID,
		CinemaID: showtime.CinemaID,
		Showtime: showtime.Showtime,
	}
	return response, nil
}

func (s *ShowtimeService) GetAllShowtime() ([]model.ShowtimeResponse, error) {
	var response []model.ShowtimeResponse
	showtimes, err := s.repository.GetAllShowtime()
	if err != nil {
		return response, err
	}
	for _, showtime := range showtimes {
		response = append(response, model.ShowtimeResponse{
			ID:       showtime.ID,
			MovieID:  showtime.MovieID,
			CinemaID: showtime.CinemaID,
			Showtime: showtime.Showtime,
		})
	}
	return response, nil
}

func (s *ShowtimeService) GetBookedSeats(id string) ([]model.BookedSeatResponse, error) {
	var response = []model.BookedSeatResponse{}
	seats, err := s.repository.GetBookedSeats(id)
	if err != nil {
		return response, err
	}
	for _, seat := range seats {
		response = append(response, model.BookedSeatResponse{
			Line:   seat.SeatLine,
			Number: seat.SeatNumber,
		})
	}
	return response, nil
}
