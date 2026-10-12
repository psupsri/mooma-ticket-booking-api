package repository

import (
	"mooma-api/model"
	"strconv"

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

func (r *ShowtimeRepository) GetAllShowtime() ([]model.Showtime, error) {
	var showtimes []model.Showtime
	result := r.db.Find(&showtimes)
	if result.Error != nil {
		return showtimes, result.Error
	}
	return showtimes, nil
}

func (r *ShowtimeRepository) GetShowtimeByMovieID(id string) ([]model.Showtime, error) {
	var showtimes []model.Showtime

	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return showtimes, err
	}

	result := r.db.Where("movie_id = ?", uintID).Find(&showtimes)
	if result.Error != nil {
		return showtimes, result.Error
	}
	return showtimes, nil
}

func (r *ShowtimeRepository) GetBookedSeats(id string) ([]model.BookedShowtimeSeat, error) {
	var seats []model.BookedShowtimeSeat

	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return seats, err
	}

	result := r.db.Where("showtime_id = ?", uintID).Find(&seats)
	if result.Error != nil {
		return seats, result.Error
	}
	return seats, nil
}
