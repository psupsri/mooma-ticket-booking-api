package repository

import (
	"mooma-api/model"

	"gorm.io/gorm"
)

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

func (r *MovieRepository) Create(movie *model.Movie) error {
	return r.db.Create(movie).Error
}

func (r *MovieRepository) GetAll() ([]model.Movie, error) {
	var movies []model.Movie
	result := r.db.Find(&movies)
	if result.Error != nil {
		return movies, result.Error
	}
	return movies, nil
}
