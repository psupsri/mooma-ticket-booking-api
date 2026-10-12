package service

import (
	"mooma-api/model"
	"mooma-api/repository"
)

type MovieService struct {
	repository         *repository.MovieRepository
	showtimeRepository *repository.ShowtimeRepository
	schedule           *Schedule
}

func NewMovieService(
	repository *repository.MovieRepository,
	showtimeRepository *repository.ShowtimeRepository,
	schedule *Schedule,
) *MovieService {
	return &MovieService{
		repository:         repository,
		showtimeRepository: showtimeRepository,
		schedule:           schedule,
	}
}

func (s *MovieService) CreateMovie(request model.CreateMovieRequest) (model.MovieResponse, error) {
	var response model.MovieResponse
	movie := model.Movie{
		Title:     request.Title,
		PosterURL: request.PosterURL,
	}
	if err := s.repository.Create(&movie); err != nil {
		return response, err
	}
	response = model.MovieResponse{
		ID:        movie.ID,
		Title:     movie.Title,
		PosterURL: movie.PosterURL,
	}
	return response, nil
}

// GetAllMovie lists movies; with nowShowing, only those with at least one
// showtime still open for booking.
func (s *MovieService) GetAllMovie(nowShowing bool) ([]model.MovieResponse, error) {
	var response = []model.MovieResponse{}
	movies, err := s.repository.GetAll()
	if err != nil {
		return response, err
	}

	var showing map[uint]bool
	if nowShowing {
		showtimes, err := s.showtimeRepository.GetAllShowtime()
		if err != nil {
			return response, err
		}
		showing = make(map[uint]bool)
		for _, showtime := range showtimes {
			if s.schedule.BookingOpen(showtime) {
				showing[showtime.MovieID] = true
			}
		}
	}

	for _, movie := range movies {
		if nowShowing && !showing[movie.ID] {
			continue
		}
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

func (s *MovieService) GetByIDWithShowtimes(id string) (model.GetMovieShowtimesResponse, error) {
	var response model.GetMovieShowtimesResponse
	movie, err := s.repository.GetByID(id)
	if err != nil {
		return response, err
	}
	showtimes, err := s.showtimeRepository.GetShowtimeByMovieID(id)
	if err != nil {
		return response, err
	}

	response = model.GetMovieShowtimesResponse{
		Movie: model.MovieResponse{
			ID:        movie.ID,
			Title:     movie.Title,
			PosterURL: movie.PosterURL,
		},
		Showtimes: s.toMovieShowtime(showtimes),
	}
	return response, err
}

func (s *MovieService) toMovieShowtime(showtimes []model.Showtime) []model.MovieShowtimeResponse {
	var response []model.MovieShowtimeResponse
	var cinemaIndex = make(map[uint]int)
	for _, showtime := range showtimes {
		index, found := cinemaIndex[showtime.CinemaID]
		if found {
			response[index].Showtimes = append(response[index].Showtimes, s.schedule.showtimeResponse(showtime))
		} else {
			cinemaIndex[showtime.CinemaID] = len(response)
			response = append(response, model.MovieShowtimeResponse{
				CinemaID:  showtime.CinemaID,
				Showtimes: []model.CinemaShowtimeResponse{s.schedule.showtimeResponse(showtime)},
			})
		}
	}
	return response
}
