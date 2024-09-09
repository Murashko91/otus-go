package storage

import (
	"context"
	"time"
)

type Storage interface {
	Connect() error
	Close() error
	CreateEvent(context.Context, Event) (Event, error)
	UpdateEvent(context.Context, Event) (Event, error)
	DeleteEvent(context.Context, int) error
	GetDailyEvents(context.Context, time.Time) ([]Event, error)
	GetWeeklyEvents(context.Context, time.Time) ([]Event, error)
	GetMonthlyEvents(context.Context, time.Time) ([]Event, error)
	GetEventsToSend(context.Context) ([]Event, error)
	DeleteOutdatedEvents(context.Context) (int, error)
}

type SenderStorage interface {
	Connect() error
	Close() error
	MarkEvent(context.Context, int) error
}

type Event struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Descr     string    `db:"descr"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	UserID    int       `db:"user_id"`
	IsSent    bool      `db:"sent"`
}

type Info struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}
