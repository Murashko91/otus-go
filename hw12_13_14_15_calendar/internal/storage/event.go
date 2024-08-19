package storage

import "time"

type Event struct {
	Id        int       `db:"id"`
	Title     string    `db:"title"`
	Descr     string    `db:"descr"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	UserId    int       `db:"user_id"`
}

type User struct {
	Id         int
	Name       string
	SecondName time.Time
	Email      string
}
