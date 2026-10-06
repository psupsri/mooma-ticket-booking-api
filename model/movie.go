package model

type Movie struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	PosterURL string
}

type CreateMovieRequest struct {
	Title     string `json:"title" binding:"required"`
	PosterURL string `json:"poster_url"`
}

type MovieResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	PosterURL string `json:"poster_url"`
}

type GetAllMovieResponse struct {
	Movies []MovieResponse `json:"movies"`
}

type GetMovieByIDResponse struct {
	Movie MovieResponse `json:"movie"`
}

type GetMovieShowtimesResponse struct {
	Movie     MovieResponse           `json:"movie"`
	Showtimes []MovieShowtimeResponse `json:"showtimes"`
}
