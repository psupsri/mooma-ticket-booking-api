package service

import (
	"mooma-api/model"
	"mooma-api/repository"
)

type MovieService struct {
	repository *repository.MovieRepository
}

func NewMovieService(repository *repository.MovieRepository) *MovieService {
	return &MovieService{
		repository: repository,
	}
}

func (s *MovieService) CreateMovie(request model.CreateMovieRequest) (model.Movie, error) {
	movie := model.Movie{
		Title:     request.Title,
		PosterURL: request.PosterURL,
	}
	if err := s.repository.Create(&movie); err != nil {
		return movie, err
	}
	return movie, nil
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

func (s *MovieService) GetByID(id string) (model.MovieResponse, error) {
	var response model.MovieResponse
	movie, err := s.repository.GetByID(id)
	if err != nil {
		return response, err
	}
	response = model.MovieResponse{
		ID:        movie.ID,
		Title:     movie.Title,
		PosterURL: movie.PosterURL,
	}
	return response, err
}
