package service

import (
	"mooma-api/model"
	"mooma-api/repository"
)

type MovieService struct {
	repository repository.MovieRepository
}

func NewMovieService(repository repository.MovieRepository) *MovieService {
	return &MovieService{
		repository: repository,
	}
}

func (s *MovieService) GetAllMovie() ([]model.MovieResponse, error) {
	var response = []model.MovieResponse{}
	movies, err := s.repository.GetAll()
	if err != nil {
		return response, err
	}

	for _, movie := range movies {
		response = append(response, model.MovieResponse{
			ID:        movie.ID,
			Title:     movie.Title,
			PosterURL: movie.PosterURL,
		})
	}

	return response, nil
}
