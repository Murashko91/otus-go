package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type memoryDb struct {
	userMap map[int]userDb
}

type userDb struct {
	storage.User
	events map[int]storage.Event
}

type Storage struct {
	mu *sync.RWMutex
	db memoryDb
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect() error {

	s.db = memoryDb{userMap: make(map[int]userDb)}
	s.mu = &sync.RWMutex{}
	return nil
}

func (s *Storage) Close() error {
	s.db = memoryDb{}
	return nil
}

var newUserId int32
var newEventId int32

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {

	userId, err := getUserIdWithCheck(ctx, event.UserId, "UpdateEvent")

	if err != nil {
		return event, err
	}
	userDb, exists := s.db.userMap[userId]

	if !exists {
		return event, fmt.Errorf("create event error: user id is missed in db: %d", userId)
	}

	eventId := int(atomic.AddInt32(&newEventId, 1))
	userDb.events[eventId] = event
	event.Id = eventId

	return event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {

	userId, err := getUserIdWithCheck(ctx, event.UserId, "UpdateEvent")

	if err != nil {
		return event, err
	}

	userDb, exists := s.db.userMap[userId]

	if !exists {
		return event, fmt.Errorf("update event error: user id is missed in db: %d", userId)
	}

	if event.UserId != userId {
		return event, fmt.Errorf("mismatch user id for update event: %d and %d", userId, event.UserId)
	}
	userDb.events[event.Id] = event

	return event, nil

}

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {

	userId, err := getUserId(ctx, "DeleteEvent")

	if err != nil {
		return err
	}
	userDb, exists := s.db.userMap[userId]

	if !exists {
		return fmt.Errorf("user id is missed in db: %d", userId)
	}

	delete(userDb.events, id)
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

	userId, err := getUserId(ctx, "getEvents")

	if err != nil {
		return nil, err
	}

	userDb, exists := s.db.userMap[userId]

	if !exists {
		return nil, fmt.Errorf("get events error: user id is missed in db: %d", userId)
	}

	result := make([]storage.Event, 0)

	for _, event := range userDb.events {
		if event.StartDate.After(startDate) && event.EndDate.Before(endDate) {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) CreateUser(ctx context.Context, user storage.User) (storage.User, error) {

	userId := int(atomic.AddInt32(&newUserId, 1))
	user.Id = userId

	userDb := userDb{User: user, events: make(map[int]storage.Event)}

	s.db.userMap[user.Id] = userDb

	return storage.User{}, nil
}

func (s *Storage) GetUser(ctx context.Context) (storage.User, error) {

	userId, err := getUserId(ctx, "GetUser")

	if err != nil {
		return storage.User{}, err
	}
	return s.db.userMap[userId].User, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user storage.User) (storage.User, error) {

	userId, err := getUserIdWithCheck(ctx, user.Id, "UpdateUser")

	if err != nil {
		return user, err
	}

	return s.db.userMap[userId].User, nil

}

func (s *Storage) DeleteUser(ctx context.Context) error {

	userId, err := getUserId(ctx, "DeleteUser")
	if err != nil {
		return err
	}
	delete(s.db.userMap, userId)

	return nil
}

func getUserIdWithCheck(ctx context.Context, id int, operationName string) (int, error) {

	userId, err := getUserId(ctx, operationName)
	if err != nil {
		return userId, err
	}

	if id != userId {
		return userId, fmt.Errorf("mismatch user id for %s: %d and %d", operationName, userId, id)
	}
	return userId, nil
}

func getUserId(ctx context.Context, operationName string) (int, error) {

	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return userId, fmt.Errorf("user id is missed in ctx for %s: %v", operationName, ctx.Value("user_id"))
	}
	return userId, nil
}

// TODO
