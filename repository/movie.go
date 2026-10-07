package repository

import (
	"errors"
	"mooma-api/model"
	"strconv"

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

func (r *MovieRepository) GetByID(id string) (model.Movie, error) {
	var movie model.Movie
	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return movie, err
	}
	result := r.db.First(&movie, uintID)
	if result.RowsAffected == 0 {
		return movie, errors.New("no movie with the id")
	}
	if result.Error != nil {
		return movie, err
	}
	return movie, err
}
