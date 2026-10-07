ALTER TABLE showtimes
ADD CONSTRAINT unique_movie_cinema_showtime
UNIQUE (movie_id, cinema_id, showtime);