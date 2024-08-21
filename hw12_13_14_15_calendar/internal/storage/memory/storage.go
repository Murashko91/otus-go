package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
)

var (
	newUserID  int32
	newEventID int32
)

type memoryDB struct {
	userMap   sync.Map
	eventsMap sync.Map
}

type Storage struct {
	db memoryDB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect() error {
	s.db = memoryDB{}
	return nil
}

func (s *Storage) Close() error {
	s.db = memoryDB{}
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	userID, err := getUserIDWithCheck(ctx, event.UserID, "CreateEvent")
	if err != nil {
		return event, err
	}
	_, exists := s.db.userMap.Load(userID)

	if !exists {
		return event, fmt.Errorf("create event error: user id is missed in db: %d", userID)
	}
	eventID := int(atomic.AddInt32(&newEventID, 1))
	event.ID = eventID
	userEventKey := getUserEventKey(userID, eventID)
	s.db.eventsMap.Store(userEventKey, event)

	return event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	userID, err := getUserIDWithCheck(ctx, event.UserID, "UpdateEvent")
	if err != nil {
		return event, err
	}

	_, exists := s.db.userMap.Load(userID)

	if !exists {
		return event, fmt.Errorf("update event error: user id is missed in db: %d", userID)
	}

	if event.UserID != userID {
		return event, fmt.Errorf("mismatch user id for update event: %d and %d", userID, event.UserID)
	}

	userEventKey := getUserEventKey(userID, event.ID)
	s.db.eventsMap.Store(userEventKey, event)

	return event, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
	userID, err := getUserID(ctx, "DeleteEvent")
	if err != nil {
		return err
	}
	_, exists := s.db.userMap.Load(userID)

	if !exists {
		return fmt.Errorf("user id is missed in db: %d", userID)
	}

	userEventKey := getUserEventKey(userID, id)
	s.db.eventsMap.Delete(userEventKey)
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

	_, exists := s.db.userMap.Load(userID)

	if !exists {
		return nil, fmt.Errorf("get events error: user id is missed in db: %d", userID)
	}

	result := make([]storage.Event, 0)

	s.db.eventsMap.Range(
		func(_, value any) bool {
			e, ok := value.(storage.Event)
			if ok &&
				e.StartDate.After(startDate) &&
				e.EndDate.Before(endDate) &&
				e.UserID == userID {
				result = append(result, e)
			}
			return true
		},
	)

	return result, nil
}

func (s *Storage) CreateUser(_ context.Context, user storage.User) (storage.User, error) {
	userID := int(atomic.AddInt32(&newUserID, 1))
	user.ID = userID

	s.db.userMap.Store(user.ID, user)

	return user, nil
}

func (s *Storage) GetUser(ctx context.Context) (storage.User, error) {
	userID, err := getUserID(ctx, "GetUser")
	if err != nil {
		return storage.User{}, err
	}
	u, ok := s.db.userMap.Load(userID)

	if !ok {
		return storage.User{}, nil
	}
	return u.(storage.User), nil
}

func (s *Storage) UpdateUser(ctx context.Context, user storage.User) (storage.User, error) {
	userID, err := getUserIDWithCheck(ctx, user.ID, "UpdateUser")
	if err != nil {
		return user, err
	}
	s.db.userMap.Store(userID, user)
	return user, nil
}

func (s *Storage) DeleteUser(ctx context.Context) error {
	userID, err := getUserID(ctx, "DeleteUser")
	if err != nil {
		return err
	}
	s.db.userMap.Delete(userID)

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

func getUserEventKey(userID, eventID int) string {
	return fmt.Sprintf("%d/%d", userID, eventID)
}
