package repository

import (
	"errors"
	"mooma-api/model"
	"strconv"

	"gorm.io/gorm"
)

type CinemaRepository struct {
	db *gorm.DB
}

func NewCinemaRepository(db *gorm.DB) *CinemaRepository {
	return &CinemaRepository{
		db: db,
	}
}

func (r *CinemaRepository) GetByID(id string) (model.Cinema, error) {
	var cinema model.Cinema
	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return cinema, err
	}
	result := r.db.First(&cinema, uintID)
	if result.RowsAffected == 0 {
		return cinema, errors.New("no cinema with the id")
	}
	if result.Error != nil {
		return cinema, result.Error
	}
	return cinema, nil
}
