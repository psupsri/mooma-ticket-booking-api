package service

import (
	"encoding/json"
	"mooma-api/model"
	"mooma-api/repository"

	"gorm.io/datatypes"
)

type CinemaService struct {
	repository *repository.CinemaRepository
}

func NewCinemaService(repository *repository.CinemaRepository) *CinemaService {
	return &CinemaService{
		repository: repository,
	}
}

func (s *CinemaService) GetByID(id string) (model.CinemaResponse, error) {
	var response model.CinemaResponse
	cinema, err := s.repository.GetByID(id)
	if err != nil {
		return response, err
	}
	seats, err := toCinemaSeatsResponse(cinema.Seats)
	if err != nil {
		return response, err
	}
	return model.CinemaResponse{
		ID:    cinema.ID,
		Name:  cinema.Name,
		Seats: seats,
	}, nil
}

func toCinemaSeatsResponse(seats datatypes.JSON) ([]model.CinemaSeatResponse, error) {
	var cinemaSeats []model.CinemaSeat
	var response = []model.CinemaSeatResponse{}
	var lineIndex = make(map[string]int)
	if err := json.Unmarshal(seats, &cinemaSeats); err != nil {
		return response, err
	}
	for _, seat := range cinemaSeats {
		index, found := lineIndex[seat.Line]
		if found {
			response[index].Numbers = append(response[index].Numbers, seat.Number)
		} else {
			lineIndex[seat.Line] = len(response)
			response = append(response, model.CinemaSeatResponse{
				Line:    seat.Line,
				Numbers: []string{seat.Number},
			})
		}
	}
	return response, nil
}
