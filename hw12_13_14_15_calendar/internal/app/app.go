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
	Debug(msg ...any)
	Info(msg ...any)
	Warn(msg ...any)
	Error(msg ...any)
}

type Storage interface {
	Connect() error
	Close() error
	CreateEvent(context.Context, storage.Event) (storage.Event, error)
	UpdateEvent(context.Context, storage.Event) (storage.Event, error)
	DeleteEvent(context.Context, int) error
	GetDailyEvents(context.Context, time.Time) ([]storage.Event, error)
	GetWeeklyEvents(context.Context, time.Time) ([]storage.Event, error)
	GetMonthlyEvents(context.Context, time.Time) ([]storage.Event, error)
	CreateUser(context.Context, storage.User) (storage.User, error)
	GetUser(context.Context) (storage.User, error)
	UpdateUser(context.Context, storage.User) (storage.User, error)
	DeleteUser(context.Context) error
}

func New(logger Logger, storage Storage) *App {
	return &App{Logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetDailyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.GetDailyEvents(ctx, date)
}

func (a *App) GetWeeklyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.GetWeeklyEvents(ctx, date)
}

func (a *App) GetMonthlyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.GetMonthlyEvents(ctx, date)
}

func (a *App) CreateUser(ctx context.Context, user storage.User) (storage.User, error) {
	return a.storage.CreateUser(ctx, user)
}
func (a *App) GetUser(ctx context.Context) (storage.User, error) {
	return a.storage.GetUser(ctx)
}
func (a *App) UpdateUser(ctx context.Context, user storage.User) (storage.User, error) {
	return a.storage.UpdateUser(ctx, user)
}
func (a *App) DeleteUser(ctx context.Context) error {
	return a.storage.DeleteUser(ctx)

}
