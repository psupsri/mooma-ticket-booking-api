package service

import (
	"time"

	"mooma-api/model"
)

const (
	// BookingGracePeriod is how long after a showtime starts it can still be booked.
	BookingGracePeriod = 15 * time.Minute
	// CancellationWindow is how long before the showtime cancelling closes.
	CancellationWindow = 2 * time.Hour
)

// Schedule applies the time rules for showtimes. Showtimes are stored as
// cinema wall-clock times (TIMESTAMP without time zone, which the driver
// labels UTC), so every rule first reads them in the cinema's zone.
type Schedule struct {
	location *time.Location
	now      func() time.Time
}

func NewSchedule(location *time.Location) *Schedule {
	return &Schedule{location: location, now: time.Now}
}

func (s *Schedule) StartsAt(showtime model.Showtime) time.Time {
	t := showtime.Showtime
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), s.location)
}

// BookingOpen reports whether seats can still be booked: until 15 minutes
// after the showtime starts.
func (s *Schedule) BookingOpen(showtime model.Showtime) bool {
	return s.now().Before(s.StartsAt(showtime).Add(BookingGracePeriod))
}

func (s *Schedule) CancelDeadline(showtime model.Showtime) time.Time {
	return s.StartsAt(showtime).Add(-CancellationWindow)
}

func (s *Schedule) CanCancel(showtime model.Showtime) bool {
	return s.now().Before(s.CancelDeadline(showtime))
}

func (s *Schedule) showtimeResponse(showtime model.Showtime) model.CinemaShowtimeResponse {
	return model.CinemaShowtimeResponse{
		ID:       showtime.ID,
		Showtime: showtime.Showtime,
		StartsAt: s.StartsAt(showtime),
		Bookable: s.BookingOpen(showtime),
	}
}
