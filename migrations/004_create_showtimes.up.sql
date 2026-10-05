CREATE TABLE showtimes (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    cinema_id INTEGER NOT NULL REFERENCES cinemas(id),
    showtime TIMESTAMP NOT NULL
);