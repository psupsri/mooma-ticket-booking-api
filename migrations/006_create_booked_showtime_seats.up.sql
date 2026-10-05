CREATE TABLE booked_showtime_seats (
    showtime_id INTEGER NOT NULL REFERENCES showtimes(id),
    booking_id INTEGER NOT NULL REFERENCES bookings(id),
    seat_line VARCHAR(10) NOT NULL,
    seat_number VARCHAR(10) NOT NULL,

    PRIMARY KEY (showtime_id, seat_line, seat_number)
);