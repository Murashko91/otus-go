package storage

import "time"

type Event struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Descr     string    `db:"descr"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	UserID    int       `db:"user_id"`
}

type User struct {
	ID    int
	Name  string
	Email string
}
