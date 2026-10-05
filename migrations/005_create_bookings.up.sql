CREATE TYPE booking_status AS ENUM (
    'confirmed',
    'cancelled',
    'expired'
);

CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL REFERENCES customers(id),
    showtime_id INTEGER NOT NULL REFERENCES showtimes(id),
    seats JSONB NOT NULL,
    status booking_status NOT NULL DEFAULT 'confirmed',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);