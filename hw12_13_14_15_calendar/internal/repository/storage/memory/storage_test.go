package memorystorage

import (
	"context"
	"testing"
	"time"

	logmock "github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/common/mocks"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/domain"
)

func TestStorage(t *testing.T) { //nolint: gocognit
	log := &logmock.Logger{}
	ctx := context.Background()
	storage := New(log)

	startDate := time.Date(2024, 0o2, 0o3, 20, 30, 50, 0, time.UTC)
	endDate := startDate.Add(time.Hour * 1)

	t.Run("create", func(t *testing.T) {
		event := domain.Event{
			UserID:       10,
			Title:        "some title",
			StartDate:    startDate,
			EndDate:      endDate,
			Notification: 0,
			Description:  "some description",
		}
		id, err := storage.Create(ctx, &event)
		if err != nil {
			t.Error(err)
		}
		if len(id) == 0 {
			t.Errorf("created id is empty")
		}
	})

	t.Run("find", func(t *testing.T) {
		conditions := []int{
			domain.AllNotifications,
			domain.TakeDayPeriodNotification,
			domain.TakeWeekPeriodNotification,
			domain.TakeMonthPeriodNotification,
		}
		filters := []domain.Event{
			{StartDate: time.Time{}, EndDate: time.Time{}},
			{
				StartDate: time.Date(2024, 0o2, 0o3, 0o0, 0o0, 0o0, 0, time.UTC),
				EndDate:   time.Date(2024, 0o2, 0o3, 23, 59, 59, 0, time.UTC),
			},
			{
				StartDate: time.Date(2024, 0o1, 29, 0o0, 0o0, 0o0, 0, time.UTC),
				EndDate:   time.Date(2024, 0o2, 0o4, 23, 59, 59, 0, time.UTC),
			},
			{
				StartDate: time.Date(2024, 0o2, 0o1, 0o0, 0o0, 0o0, 0, time.UTC),
				EndDate:   time.Date(2024, 0o2, 29, 23, 59, 59, 0, time.UTC),
			},
		}
		for i := range conditions {
			events, err := storage.Find(ctx, conditions[i], filters[i])
			if err != nil {
				t.Error(err)
			}
			if len(events) != 1 {
				t.Errorf("condition=%d, wait len of events=%d, but got=%d", conditions[i], 1, len(events))
			}
		}
	})

	t.Run("update", func(t *testing.T) {
		events, err := storage.Find(ctx, domain.AllNotifications, domain.Event{})
		if err != nil {
			t.Error(err)
		}
		if len(events) != 1 {
			t.Errorf("condition=%d, wait len of events=%d, but got=%d", domain.AllNotifications, 1, len(events))
		}

		event := events[0]
		event.Title = "new title"
		event.UserID = 50
		err = storage.Update(ctx, &event)
		if err != nil {
			t.Error(err)
		}

		events, err = storage.Find(ctx, domain.AllNotifications, domain.Event{})
		if err != nil {
			t.Error(err)
		}
		if len(events) != 1 {
			t.Errorf("condition=%d, wait len of events=%d, but got=%d", domain.AllNotifications, 1, len(events))
		}
		if events[0].UserID != 50 && events[0].Title != "new title" {
			t.Errorf("operation update was bad")
		}
	})

	t.Run("delete", func(t *testing.T) {
		events, err := storage.Find(ctx, domain.AllNotifications, domain.Event{})
		if err != nil {
			t.Error(err)
		}
		if len(events) != 1 {
			t.Errorf("condition=%d, wait len of events=%d, but got=%d", domain.AllNotifications, 1, len(events))
		}

		err = storage.Delete(ctx, events[0].ID)
		if err != nil {
			t.Error(err)
		}

		events, err = storage.Find(ctx, domain.AllNotifications, domain.Event{})
		if err != nil {
			t.Error(err)
		}
		if len(events) != 0 {
			t.Errorf("operation delete was bad")
		}
	})
}
