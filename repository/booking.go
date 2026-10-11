package repository

import (
	"mooma-api/model"

	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

func (r *BookingRepository) CreateBooking(booking *model.Booking, seats []model.CinemaSeat) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create booking
		if err := tx.Create(booking).Error; err != nil {
			return err
		}

		// Create booked seats
		for _, seat := range seats {
			bookedShowtimeSeat := model.BookedShowtimeSeat{
				ShowtimeID: booking.ShowtimeID,
				BookingID:  booking.ID,
				SeatLine:   seat.Line,
				SeatNumber: seat.Number,
			}

			if err := tx.Create(&bookedShowtimeSeat).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
