package repository

import (
	"mooma-api/model"

	"gorm.io/gorm"
)

type ShowtimeRepository struct {
	db *gorm.DB
}

func NewShowtimeRepository(db *gorm.DB) *ShowtimeRepository {
	return &ShowtimeRepository{
		db: db,
	}
}

func (r *ShowtimeRepository) CreateShowtime(request *model.Showtime) error {
	return r.db.Create(request).Error
}
