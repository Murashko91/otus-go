package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestEventStorage(t *testing.T) {
	count := 100
	newName := "TestUpdated"

	t.Run("test event db CRUD", func(t *testing.T) {
		memory := New()
		userID := 1

		ctx := app.SetContextValue(context.Background(), app.UserIDKey, userID)

		wg := &sync.WaitGroup{}
		wg.Add(count)

		// create Events
		for i := 0; i < count; i++ {
			currentTime := time.Now()
			go func() {
				newEvent, err := memory.CreateEvent(ctx,
					storage.Event{
						Title:     "Test",
						Descr:     "test",
						StartDate: currentTime,
						EndDate:   currentTime.Add(time.Hour * 24 * time.Duration(i)),
						UserID:    userID,
					})
				wg.Done()
				if err != nil {
					require.Nilf(t, err, err.Error())
				}
				require.NotNilf(t, newEvent.ID, fmt.Sprintf("id should be created for newly created event record: %d", i))
			}()
		}
		wg.Wait()

		// update Events
		wg.Add(count)
		for i := 0; i < count; i++ {
			currentTime := time.Now()
			go func() {
				newEvent, err := memory.UpdateEvent(ctx,
					storage.Event{
						Title:     newName,
						Descr:     "test",
						StartDate: currentTime.Add(time.Second),
						EndDate:   currentTime.Add(time.Hour * 24 * time.Duration(i)),
						UserID:    userID,
						ID:        i,
					})
				wg.Done()
				if err != nil {
					require.Nilf(t, err, err.Error())
				}
				require.NotNilf(t, newEvent.ID, fmt.Sprintf("id should be created for newly created event record: %d", i))
			}()
		}
		wg.Wait()

		// test get events
		wg.Add(3)
		eventsCountMap := map[int]int{
			0: 2,
			1: 8,
			2: 31,
		}
		for i := 0; i < 3; i++ {
			go func() {
				getEventsFunc := getEventQueryMetod(i, memory)
				events, err := getEventsFunc(ctx, time.Now())
				require.Equal(t, len(events), eventsCountMap[i], "Incorrect events count")

				if err != nil {
					require.Nilf(t, err, err.Error())
				}
				wg.Done()
			}()
		}
		wg.Wait()

		// delete Events
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func() {
				err := memory.DeleteEvent(ctx, i)

				if err != nil {
					require.Nilf(t, err, err.Error())
				}
				wg.Done()
			}()
		}

		wg.Wait()
		events, err := memory.GetMonthlyEvents(ctx, time.Now().Add(time.Hour*-1))
		if err != nil {
			require.Nilf(t, err, err.Error())
		}

		require.Equal(t, len(events), 0, "events has not been deleted")
	})
}

func getEventQueryMetod(i int, memory *Storage) func(ctx context.Context,
	startDate time.Time) ([]storage.Event, error) {
	getEventsFunc := memory.GetDailyEvents

	switch i {
	case 1:
		getEventsFunc = memory.GetWeeklyEvents
	case 2:
		getEventsFunc = memory.GetMonthlyEvents
	}
	return getEventsFunc
}
