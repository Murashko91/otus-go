package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	// TODO
	mu sync.RWMutex //nolint:unused
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) CreateEvent(event storage.Event) (storage.Event, error) {
	// TODO
	return storage.Event{}, nil
}

func (s *Storage) UpdateEvent(id string, event storage.Event) (storage.Event, error) {
	return storage.Event{}, nil
}

func (s *Storage) DeleteEvent(id string) (storage.Event, error) {
	// TODO
	return storage.Event{}, nil
}

func (s *Storage) GetDailyEvents(startDate time.Time) ([]storage.Event, error) {
	return nil, nil
}

func (s *Storage) GetMonthlyEvents(startDate time.Time) ([]storage.Event, error) {
	return nil, nil
}

func (s *Storage) GetWeeklyEvents(startDate time.Time) ([]storage.Event, error) {
	return nil, nil
}

// TODO
