package service

import (
	"mooma-api/model"
	"testing"
	"time"
)

func newTestSchedule(t *testing.T, now time.Time) *Schedule {
	t.Helper()
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		t.Fatal(err)
	}
	s := NewSchedule(bangkok)
	s.now = func() time.Time { return now }
	return s
}

// Showtimes come back from Postgres as wall-clock times labelled UTC.
func wallClock(value string) model.Showtime {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return model.Showtime{Showtime: t}
}

func TestStartsAtUsesCinemaTimezone(t *testing.T) {
	s := newTestSchedule(t, time.Now())

	got := s.StartsAt(wallClock("2026-10-07T10:30:00Z"))

	// 10:30 in Bangkok (UTC+7) is 03:30 UTC.
	want := time.Date(2026, 10, 7, 3, 30, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("startsAt = %v, want %v", got, want)
	}
}

func TestCanCancel(t *testing.T) {
	showtime := wallClock("2026-10-07T19:00:00Z") // 19:00 Bangkok = 12:00 UTC
	tests := []struct {
		name string
		now  time.Time
		want bool
	}{
		{"a day before", time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC), true},
		{"just over 2 hours before", time.Date(2026, 10, 7, 9, 59, 59, 0, time.UTC), true},
		{"exactly 2 hours before", time.Date(2026, 10, 7, 10, 0, 0, 0, time.UTC), false},
		{"1 hour before", time.Date(2026, 10, 7, 11, 0, 0, 0, time.UTC), false},
		{"after it started", time.Date(2026, 10, 7, 13, 0, 0, 0, time.UTC), false},
		// 21:00 in Bangkok, after the start; reading 19:00 as UTC would put the
		// deadline at 17:00 UTC and wrongly allow this.
		{"wall-clock read as UTC would wrongly allow", time.Date(2026, 10, 7, 14, 0, 0, 0, time.UTC), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestSchedule(t, tt.now)
			if got := s.CanCancel(showtime); got != tt.want {
				t.Fatalf("canCancel at %v = %v, want %v (deadline %v)", tt.now, got, tt.want, s.CancelDeadline(showtime))
			}
		})
	}
}

func TestBookingOpen(t *testing.T) {
	showtime := wallClock("2026-10-07T19:00:00Z") // 19:00 Bangkok = 12:00 UTC
	tests := []struct {
		name string
		now  time.Time
		want bool
	}{
		{"before it starts", time.Date(2026, 10, 7, 11, 0, 0, 0, time.UTC), true},
		{"at the start", time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC), true},
		{"14 minutes in", time.Date(2026, 10, 7, 12, 14, 59, 0, time.UTC), true},
		{"15 minutes in", time.Date(2026, 10, 7, 12, 15, 0, 0, time.UTC), false},
		{"an hour in", time.Date(2026, 10, 7, 13, 0, 0, 0, time.UTC), false},
		// 18:00 UTC is before 19:00 if the zone were ignored, but it is 01:00
		// the next day in Bangkok.
		{"wall-clock read as UTC would wrongly allow", time.Date(2026, 10, 7, 18, 0, 0, 0, time.UTC), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestSchedule(t, tt.now)
			if got := s.BookingOpen(showtime); got != tt.want {
				t.Fatalf("BookingOpen at %v = %v, want %v", tt.now, got, tt.want)
			}
		})
	}
}
