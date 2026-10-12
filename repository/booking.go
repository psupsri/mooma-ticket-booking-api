package repository

import (
	"errors"
	"mooma-api/model"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrSeatAlreadyBooked = errors.New("one or more seats are already booked")
	ErrBookingNotFound   = errors.New("booking not found")
	ErrShowtimeNotFound  = errors.New("showtime not found")
)

const uniqueViolation = "23505"

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

// CreateBooking stores the booking and its seats. canBook runs on the
// showtime first, so a closed showtime is refused before anything is written.
func (r *BookingRepository) CreateBooking(booking *model.Booking, seats []model.CinemaSeat, canBook func(model.Showtime) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var showtime model.Showtime
		result := tx.Clauses(clause.Locking{Strength: "SHARE"}).Limit(1).Find(&showtime, booking.ShowtimeID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrShowtimeNotFound
		}
		if err := canBook(showtime); err != nil {
			return err
		}

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
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
					return ErrSeatAlreadyBooked
				}
				return err
			}
		}

		return nil
	})
}

func (r *BookingRepository) GetByIDWithDetails(id uint) (model.Booking, model.Movie, error) {
	var booking model.Booking
	var movie model.Movie
	if err := r.db.Preload("Customer").Preload("Showtime").First(&booking, id).Error; err != nil {
		return booking, movie, err
	}
	if err := r.db.First(&movie, booking.Showtime.MovieID).Error; err != nil {
		return booking, movie, err
	}
	return booking, movie, nil
}

// Bookings of a customer, newest first, with the movie of each showtime keyed by ID.
func (r *BookingRepository) GetByCustomerID(customerID uint) ([]model.Booking, map[uint]model.Movie, error) {
	var bookings []model.Booking
	movies := make(map[uint]model.Movie)
	if err := r.db.Preload("Customer").Preload("Showtime").
		Where("customer_id = ?", customerID).
		Order("created_at DESC, id DESC").
		Find(&bookings).Error; err != nil {
		return bookings, movies, err
	}

	var movieIDs []uint
	for _, booking := range bookings {
		movieIDs = append(movieIDs, booking.Showtime.MovieID)
	}
	if len(movieIDs) == 0 {
		return bookings, movies, nil
	}

	var found []model.Movie
	if err := r.db.Where("id IN ?", movieIDs).Find(&found).Error; err != nil {
		return bookings, movies, err
	}
	for _, movie := range found {
		movies[movie.ID] = movie
	}
	return bookings, movies, nil
}

// Delete removes a customer's booking and releases its seats. A booking that
// doesn't exist or belongs to another customer is reported as not found.
// canDelete runs on the locked booking before anything is removed.
func (r *BookingRepository) Delete(id uint, customerID uint, canDelete func(model.Booking) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var booking model.Booking
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Showtime").
			Where("id = ? AND customer_id = ?", id, customerID).Limit(1).Find(&booking)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrBookingNotFound
		}
		if err := canDelete(booking); err != nil {
			return err
		}

		// Seats reference the booking, so they go first.
		if err := tx.Where("booking_id = ?", booking.ID).Delete(&model.BookedShowtimeSeat{}).Error; err != nil {
			return err
		}
		return tx.Delete(&booking).Error
	})
}
