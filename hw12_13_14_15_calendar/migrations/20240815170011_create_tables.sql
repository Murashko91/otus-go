-- +goose Up
-- +goose StatementBegin

CREATE DATABASE calendar;


CREATE TABLE IF NOT EXISTS events (
  id SERIAL PRIMARY KEY,
  user_id INT,
  title TEXT NOT NULL,
  descr TEXT,
  start_date  TIMESTAMP NOT NULL,
  end_date  TIMESTAMP NOT NULL
	);

--INSERT INTO events (title, event_date, user_id) VALUES ( 'Some event', '2024-08-14 16:50:36', '1');
--select * from events;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE events;
-- +goose StatementEnd
