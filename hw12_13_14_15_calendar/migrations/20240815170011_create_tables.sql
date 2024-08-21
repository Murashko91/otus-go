-- +goose Up
-- +goose StatementBegin

CREATE DATABASE calendar;


CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY ,
  name TEXT NOT NULL,
  email TEXT NOT NULL
);

--INSERT INTO users (name, email) VALUES ( 'Sergey Murashko', 'sergey-test@test.ru');
--select * from users;


CREATE TABLE IF NOT EXISTS events (
  id SERIAL PRIMARY KEY,
  user_id INT,
  title TEXT NOT NULL,
  event_date  TIMESTAMP NOT NULL,
	last_modified_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_user
    FOREIGN KEY(user_id) 
	    REFERENCES users(id)
	     ON DELETE CASCADE
	);

--INSERT INTO events (title, event_date, user_id) VALUES ( 'Some event', '2024-08-14 16:50:36', '1');
--select * from events;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE users;
DROP TABLE events;
-- +goose StatementEnd
