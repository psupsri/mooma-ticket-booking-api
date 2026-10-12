package main

// API tests: real router, real handlers/services/repositories, real Postgres.
// The booking rules depend on SQL (row locks, the unique seat key, two-table
// deletes inside transactions), so they are tested against a database rather
// than mocks.
//
// They run only when MOOMA_TEST_DATABASE_DSN points at a database whose name
// ends in "_test". Each test drops and recreates its schema, so never point it
// at real data. Example:
//
//	MOOMA_TEST_DATABASE_DSN="host=localhost port=5432 user=... password=... dbname=mooma_test sslmode=disable" go test ./...

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"mooma-api/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const testCinemaZone = "Asia/Bangkok"

type api struct {
	t      *testing.T
	router *gin.Engine
	db     *gorm.DB
	zone   *time.Location
}

func newAPI(t *testing.T) *api {
	t.Helper()
	dsn := os.Getenv("MOOMA_TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("MOOMA_TEST_DATABASE_DSN not set; skipping API tests")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatalf("connect test database: %v", err)
	}
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	})

	var name string
	db.Raw("SELECT current_database()").Scan(&name)
	if !strings.HasSuffix(name, "_test") {
		t.Fatalf("refusing to reset database %q: its name must end in _test", name)
	}
	resetSchema(t, db)

	zone, err := time.LoadLocation(testCinemaZone)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("JWT_SECRET", "api-test-secret")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	registerApp(r, db, service.NewSchedule(zone))

	return &api{t: t, router: r, db: db, zone: zone}
}

// resetSchema recreates the schema from the project's own migrations.
func resetSchema(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public").Error; err != nil {
		t.Fatalf("reset schema: %v", err)
	}
	files, err := filepath.Glob("migrations/*.up.sql")
	if err != nil || len(files) == 0 {
		t.Fatalf("find migrations: %v", err)
	}
	sort.Strings(files)
	for _, file := range files {
		sql, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if err := db.Exec(string(sql)).Error; err != nil {
			t.Fatalf("apply %s: %v", file, err)
		}
	}
}

// --- HTTP helpers ---

type response struct {
	status int
	body   map[string]any
}

func (a *api) do(method, path, token string, body any) response {
	a.t.Helper()
	var reader *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			a.t.Fatal(err)
		}
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	a.router.ServeHTTP(rec, req)

	var parsed map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &parsed)
	return response{status: rec.Code, body: parsed}
}

func (a *api) expect(res response, status int) response {
	a.t.Helper()
	if res.status != status {
		a.t.Fatalf("status = %d, want %d; body: %v", res.status, status, res.body)
	}
	return res
}

// --- fixtures ---

// customer registers through the API and returns a login token.
func (a *api) customer(email string) string {
	a.t.Helper()
	a.expect(a.do(http.MethodPost, "/auth/register", "", map[string]string{
		"first_name": "Test", "last_name": "User", "email": email, "password": "secret-password",
	}), http.StatusCreated)
	res := a.expect(a.do(http.MethodPost, "/auth/login", "", map[string]string{
		"email": email, "password": "secret-password",
	}), http.StatusOK)
	return res.body["AccessToken"].(string)
}

func (a *api) movie(title string) int {
	a.t.Helper()
	var id int
	a.db.Raw("INSERT INTO movies (title) VALUES (?) RETURNING id", title).Scan(&id)
	return id
}

// cinema has rows a and b with seats 1-3.
func (a *api) cinema() int {
	a.t.Helper()
	var seats []map[string]string
	for _, line := range []string{"a", "b"} {
		for _, number := range []string{"1", "2", "3"} {
			seats = append(seats, map[string]string{"line": line, "number": number})
		}
	}
	raw, _ := json.Marshal(seats)
	var id int
	a.db.Raw("INSERT INTO cinemas (name, seats) VALUES ('A', ?) RETURNING id", string(raw)).Scan(&id)
	return id
}

// showtime stores a showtime `fromNow` away, as cinema wall-clock time the
// way real showtimes are stored.
func (a *api) showtime(movieID, cinemaID int, fromNow time.Duration) int {
	a.t.Helper()
	wallClock := time.Now().In(a.zone).Add(fromNow).Truncate(time.Second).Format("2006-01-02 15:04:05")
	var id int
	a.db.Raw("INSERT INTO showtimes (movie_id, cinema_id, showtime) VALUES (?, ?, ?) RETURNING id",
		movieID, cinemaID, wallClock).Scan(&id)
	return id
}

func (a *api) count(query string, args ...any) int {
	a.t.Helper()
	var n int
	if err := a.db.Raw(query, args...).Scan(&n).Error; err != nil {
		a.t.Fatal(err)
	}
	return n
}

func seats(keys ...string) []map[string]string {
	var out []map[string]string
	for _, key := range keys {
		out = append(out, map[string]string{"Line": key[:1], "Number": key[1:]})
	}
	return out
}

func booking(showtimeID int, keys ...string) map[string]any {
	return map[string]any{"showtime_id": showtimeID, "seats": seats(keys...)}
}

func (a *api) bookedSeats(showtimeID int) []string {
	a.t.Helper()
	res := a.expect(a.do(http.MethodGet, fmt.Sprintf("/showtime/%d/booked-seats", showtimeID), "", nil), http.StatusOK)
	var out []string
	for _, s := range res.body["data"].([]any) {
		seat := s.(map[string]any)
		out = append(out, seat["line"].(string)+seat["number"].(string))
	}
	sort.Strings(out)
	return out
}

type world struct {
	*api
	alice, bob string
	movieID    int
	open       int // tomorrow: bookable and cancellable
	soon       int // in 1 hour: bookable, cancellation closed
	closed     int // started 20 minutes ago: booking closed
}

func newWorld(t *testing.T) *world {
	a := newAPI(t)
	movieID := a.movie("Titanium")
	cinemaID := a.cinema()
	return &world{
		api:     a,
		alice:   a.customer("alice@mooma.test"),
		bob:     a.customer("bob@mooma.test"),
		movieID: movieID,
		open:    a.showtime(movieID, cinemaID, 24*time.Hour),
		soon:    a.showtime(movieID, cinemaID, time.Hour),
		closed:  a.showtime(movieID, cinemaID, -20*time.Minute),
	}
}

func bookingID(res response) int {
	return int(res.body["message"].(map[string]any)["id"].(float64))
}

// --- tests ---

func TestBookingEndpointsRequireAuth(t *testing.T) {
	w := newWorld(t)

	w.expect(w.do(http.MethodPost, "/booking/", "", booking(w.open, "a1")), http.StatusUnauthorized)
	w.expect(w.do(http.MethodGet, "/booking/", "", nil), http.StatusUnauthorized)
	w.expect(w.do(http.MethodDelete, "/booking/1", "", nil), http.StatusUnauthorized)

	if n := w.count("SELECT count(*) FROM bookings"); n != 0 {
		t.Fatalf("bookings = %d after unauthenticated requests, want 0", n)
	}
}

func TestCreateBooking(t *testing.T) {
	w := newWorld(t)

	res := w.expect(w.do(http.MethodPost, "/booking/", w.alice, booking(w.open, "a1", "a2")), http.StatusOK)

	got := res.body["message"].(map[string]any)
	if got["movie"].(map[string]any)["title"] != "Titanium" {
		t.Errorf("movie = %v, want Titanium", got["movie"])
	}
	if got["customer"].(map[string]any)["email"] != "alice@mooma.test" {
		t.Errorf("customer = %v, want alice", got["customer"])
	}
	if got["can_cancel"] != true {
		t.Errorf("can_cancel = %v for tomorrow's showtime, want true", got["can_cancel"])
	}
	if seats := w.bookedSeats(w.open); strings.Join(seats, ",") != "a1,a2" {
		t.Errorf("booked seats = %v, want [a1 a2]", seats)
	}
}

func TestCreateBookingRejectsTakenSeatAndRollsBack(t *testing.T) {
	w := newWorld(t)
	w.expect(w.do(http.MethodPost, "/booking/", w.alice, booking(w.open, "a1")), http.StatusOK)

	// a2 is free but a1 is taken: nothing of Bob's request may be kept.
	w.expect(w.do(http.MethodPost, "/booking/", w.bob, booking(w.open, "a2", "a1")), http.StatusConflict)

	if n := w.count("SELECT count(*) FROM bookings"); n != 1 {
		t.Fatalf("bookings = %d, want 1 (the failed booking must roll back)", n)
	}
	if seats := w.bookedSeats(w.open); strings.Join(seats, ",") != "a1" {
		t.Fatalf("booked seats = %v, want only [a1]", seats)
	}
	w.expect(w.do(http.MethodPost, "/booking/", w.bob, booking(w.open, "a2")), http.StatusOK)
}

func TestCreateBookingValidation(t *testing.T) {
	w := newWorld(t)

	tests := []struct {
		name   string
		body   map[string]any
		status int
	}{
		{"no seats", booking(w.open), http.StatusBadRequest},
		{"unknown showtime", booking(999999, "a1"), http.StatusNotFound},
		{"started more than 15 minutes ago", booking(w.closed, "a1"), http.StatusConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w.expect(w.do(http.MethodPost, "/booking/", w.alice, tt.body), tt.status)
		})
	}
	if n := w.count("SELECT count(*) FROM bookings"); n != 0 {
		t.Fatalf("bookings = %d after rejected requests, want 0", n)
	}
}

func TestBookingHistoryIsPerCustomer(t *testing.T) {
	w := newWorld(t)
	w.expect(w.do(http.MethodPost, "/booking/", w.alice, booking(w.open, "a1")), http.StatusOK)
	w.expect(w.do(http.MethodPost, "/booking/", w.alice, booking(w.soon, "b1")), http.StatusOK)

	alice := w.expect(w.do(http.MethodGet, "/booking/", w.alice, nil), http.StatusOK).body["data"].([]any)
	bob := w.expect(w.do(http.MethodGet, "/booking/", w.bob, nil), http.StatusOK).body["data"].([]any)

	if len(alice) != 2 {
		t.Fatalf("alice sees %d bookings, want 2", len(alice))
	}
	if len(bob) != 0 {
		t.Fatalf("bob sees %d bookings, want 0: %v", len(bob), bob)
	}
	canCancel := map[float64]any{}
	for _, b := range alice {
		item := b.(map[string]any)
		canCancel[item["showtime"].(map[string]any)["id"].(float64)] = item["can_cancel"]
	}
	if canCancel[float64(w.open)] != true || canCancel[float64(w.soon)] != false {
		t.Fatalf("can_cancel by showtime = %v, want open=true soon=false", canCancel)
	}
}

func TestCancelBookingOnlyByOwner(t *testing.T) {
	w := newWorld(t)
	id := bookingID(w.expect(w.do(http.MethodPost, "/booking/", w.alice, booking(w.open, "a1", "a2")), http.StatusOK))
	path := fmt.Sprintf("/booking/%d", id)

	// Bob gets 404, not 403, so booking IDs can't be probed.
	w.expect(w.do(http.MethodDelete, path, w.bob, nil), http.StatusNotFound)
	if n := w.count("SELECT count(*) FROM booked_showtime_seats WHERE booking_id = ?", id); n != 2 {
		t.Fatalf("seats after bob's attempt = %d, want 2", n)
	}

	w.expect(w.do(http.MethodDelete, path, w.alice, nil), http.StatusOK)

	if n := w.count("SELECT count(*) FROM bookings WHERE id = ?", id); n != 0 {
		t.Errorf("booking row still exists after cancel")
	}
	if n := w.count("SELECT count(*) FROM booked_showtime_seats WHERE booking_id = ?", id); n != 0 {
		t.Errorf("booked seats = %d after cancel, want 0", n)
	}
	w.expect(w.do(http.MethodDelete, path, w.alice, nil), http.StatusNotFound)
	// The released seats can be booked again.
	w.expect(w.do(http.MethodPost, "/booking/", w.bob, booking(w.open, "a1", "a2")), http.StatusOK)
}

func TestCancelBookingClosesTwoHoursBeforeShowtime(t *testing.T) {
	w := newWorld(t)
	id := bookingID(w.expect(w.do(http.MethodPost, "/booking/", w.alice, booking(w.soon, "a1")), http.StatusOK))

	w.expect(w.do(http.MethodDelete, fmt.Sprintf("/booking/%d", id), w.alice, nil), http.StatusConflict)

	if n := w.count("SELECT count(*) FROM booked_showtime_seats WHERE booking_id = ?", id); n != 1 {
		t.Fatalf("seats = %d after refused cancel, want 1", n)
	}
}

func TestNowShowingListsOnlyBookableMovies(t *testing.T) {
	a := newAPI(t)
	cinemaID := a.cinema()
	showing := a.movie("Showing")
	a.showtime(showing, cinemaID, 10*time.Minute)
	finished := a.movie("Finished")
	a.showtime(finished, cinemaID, -time.Hour)
	a.movie("No showtimes")

	titles := func(path string) []string {
		var out []string
		for _, m := range a.expect(a.do(http.MethodGet, path, "", nil), http.StatusOK).body["data"].([]any) {
			out = append(out, m.(map[string]any)["title"].(string))
		}
		sort.Strings(out)
		return out
	}

	if got := titles("/movie/?now_showing=true"); strings.Join(got, ",") != "Showing" {
		t.Errorf("now showing = %v, want [Showing]", got)
	}
	if got := titles("/movie/"); len(got) != 3 {
		t.Errorf("all movies = %v, want 3", got)
	}
}
