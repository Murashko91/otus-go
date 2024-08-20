package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type memoryDB struct {
	userMap map[int]userDB
}

type userDB struct {
	storage.User
	events map[int]storage.Event
}

type Storage struct {
	mu *sync.RWMutex
	db memoryDB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect() error {
	s.db = memoryDB{userMap: make(map[int]userDB)}
	s.mu = &sync.RWMutex{}
	return nil
}

func (s *Storage) Close() error {
	s.db = memoryDB{}
	return nil
}

var (
	newUserID  int32
	newEventID int32
)

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	userID, err := getUserIDWithCheck(ctx, event.UserID, "UpdateEvent")
	if err != nil {
		return event, err
	}
	userDB, exists := s.db.userMap[userID]

	if !exists {
		return event, fmt.Errorf("create event error: user id is missed in db: %d", userID)
	}

	eventID := int(atomic.AddInt32(&newEventID, 1))
	userDB.events[eventID] = event
	event.ID = eventID

	return event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	userID, err := getUserIDWithCheck(ctx, event.UserID, "UpdateEvent")
	if err != nil {
		return event, err
	}

	userDB, exists := s.db.userMap[userID]

	if !exists {
		return event, fmt.Errorf("update event error: user id is missed in db: %d", userID)
	}

	if event.UserID != userID {
		return event, fmt.Errorf("mismatch user id for update event: %d and %d", userID, event.UserID)
	}
	s.mu.Lock()
	userDB.events[event.ID] = event
	s.mu.Unlock()

	return event, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
	userID, err := getUserID(ctx, "DeleteEvent")
	if err != nil {
		return err
	}
	userDB, exists := s.db.userMap[userID]

	if !exists {
		return fmt.Errorf("user id is missed in db: %d", userID)
	}

	delete(userDB.events, id)
	return nil
}

func (s *Storage) GetDailyEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	endDate := startDate.Add(time.Hour * 24)

	return s.getEvents(ctx, startDate, endDate)
}

func (s *Storage) GetMonthlyEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	endDate := startDate.Add(time.Hour * 24 * 30)

	return s.getEvents(ctx, startDate, endDate)
}

func (s *Storage) GetWeeklyEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	endDate := startDate.Add(time.Hour * 24 * 7)

	return s.getEvents(ctx, startDate, endDate)
}

func (s *Storage) getEvents(ctx context.Context, startDate time.Time, endDate time.Time) ([]storage.Event, error) {
	userID, err := getUserID(ctx, "getEvents")
	if err != nil {
		return nil, err
	}

	userDB, exists := s.db.userMap[userID]

	if !exists {
		return nil, fmt.Errorf("get events error: user id is missed in db: %d", userID)
	}

	result := make([]storage.Event, 0)

	for _, event := range userDB.events {
		if event.StartDate.After(startDate) && event.EndDate.Before(endDate) {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) CreateUser(ctx context.Context, user storage.User) (storage.User, error) {
	userID := int(atomic.AddInt32(&newUserID, 1))
	user.ID = userID

	userDB := userDB{User: user, events: make(map[int]storage.Event)}

	s.db.userMap[user.ID] = userDB

	return storage.User{}, nil
}

func (s *Storage) GetUser(ctx context.Context) (storage.User, error) {
	userID, err := getUserID(ctx, "GetUser")
	if err != nil {
		return storage.User{}, err
	}
	return s.db.userMap[userID].User, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user storage.User) (storage.User, error) {
	userID, err := getUserIDWithCheck(ctx, user.ID, "UpdateUser")
	if err != nil {
		return user, err
	}

	userDB := s.db.userMap[userID]
	userDB.User = user
	s.mu.Lock()
	s.db.userMap[userID] = userDB
	s.mu.Unlock()
	return s.db.userMap[userID].User, nil
}

func (s *Storage) DeleteUser(ctx context.Context) error {
	userID, err := getUserID(ctx, "DeleteUser")
	if err != nil {
		return err
	}
	delete(s.db.userMap, userID)

	return nil
}

func getUserIDWithCheck(ctx context.Context, id int, operationName string) (int, error) {
	userID, err := getUserID(ctx, operationName)
	if err != nil {
		return userID, err
	}

	if id != userID {
		return userID, fmt.Errorf("mismatch user id for %s: %d and %d", operationName, userID, id)
	}
	return userID, nil
}

func getUserID(ctx context.Context, operationName string) (int, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return userID, fmt.Errorf("user id is missed in ctx for %s: %v", operationName, ctx.Value("user_id"))
	}
	return userID, nil
}
