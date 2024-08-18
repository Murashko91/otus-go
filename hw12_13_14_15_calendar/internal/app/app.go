package app

import (
	"context"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage Storage
	Logger
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	Connect() error
	Close() error
	CreateEvent(storage.Event) (storage.Event, error)
	UpdateEvent(string, storage.Event) (storage.Event, error)
	DeleteEvent(string) (storage.Event, error)
	GetDailyEvents(time.Time) ([]storage.Event, error)
	GetWeeklyEvents(time.Time) ([]storage.Event, error)
	GetMonthlyEvents(time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{Logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) (storage.Event, error) {
	return a.storage.CreateEvent(storage.Event{Id: id, Title: title})
}
