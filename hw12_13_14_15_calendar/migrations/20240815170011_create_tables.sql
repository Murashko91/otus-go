-- +goose Up
-- +goose StatementBegin

CREATE DATABASE calendar;


CREATE TABLE IF NOT EXISTS events (
  id SERIAL PRIMARY KEY,
  user_id INT,
  title TEXT NOT NULL,
  event_date  TIMESTAMP NOT NULL,
	last_modified_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

--INSERT INTO events (title, event_date, user_id) VALUES ( 'Some event', '2024-08-14 16:50:36', '1');
--select * from events;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE events;
-- +goose StatementEnd
