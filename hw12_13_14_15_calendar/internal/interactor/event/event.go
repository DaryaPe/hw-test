package event

import (
	"context"
	"time"

	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/common"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/domain"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/pkg/util"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/pkg/errors"
)

type Inter struct {
	logger common.Logger
	db     Storage
}

type Storage interface {
	Create(ctx context.Context, event *domain.Event) (string, error)
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, condition int, filter domain.Event) ([]domain.Event, error)
}

func New(storage Storage) *Inter {
	return &Inter{
		db: storage,
	}
}

// Create создаёт уведомление.
func (inter *Inter) Create(ctx context.Context, event *domain.Event) (string, error) {
	err := inter.checkTimeIsFree(ctx, event)
	if err != nil {
		return "", errors.Wrap(err, "inter.checkTimeIsFree")
	}

	id, err := inter.db.Create(ctx, event)
	if err != nil {
		return "", errors.Wrap(err, "inter.db.Create")
	}
	return id, nil
}

// Update обновляет уведомление.
func (inter *Inter) Update(ctx context.Context, event *domain.Event) error {
	err := inter.checkTimeIsFree(ctx, event)
	if err != nil {
		return errors.Wrap(err, "inter.checkTimeIsFree")
	}

	err = inter.checkEventExists(ctx, event.ID)
	if err != nil {
		return errors.Wrap(err, "inter.checkEventExists")
	}

	if err := inter.db.Update(ctx, event); err != nil {
		return errors.Wrap(err, "inter.db.Update")
	}
	return nil
}

// Delete удаляет уведомление.
func (inter *Inter) Delete(ctx context.Context, id string) error {
	err := inter.checkEventExists(ctx, id)
	if err != nil {
		return errors.Wrap(err, "inter.checkEventExists")
	}

	if err = inter.db.Delete(ctx, id); err != nil {
		return errors.Wrap(err, "inter.db.Delete")
	}
	return nil
}

// Find возвращает уведомления по условию.
func (inter *Inter) Find(ctx context.Context, date time.Time, condition int) ([]domain.Event, error) {
	if date.IsZero() {
		condition = domain.AllNotifications
	}

	var filter domain.Event
	switch condition {
	case domain.AllNotifications:
		filter.StartDate, filter.EndDate = time.Time{}, time.Time{}
	case domain.TakeDayPeriodNotification:
		filter.StartDate, filter.EndDate = util.StartDateDay(date), util.EndDateDay(date)
	case domain.TakeWeekPeriodNotification:
		filter.StartDate, filter.EndDate = util.StartDateWeek(date), util.EndDateWeek(date)
	case domain.TakeMonthPeriodNotification:
		filter.StartDate, filter.EndDate = util.StartDateMonth(date), util.EndDateMonth(date)
	default:
		return nil, domain.ErrNotDefinedPeriod
	}

	events, err := inter.db.Find(ctx, condition, filter)
	if err != nil {
		return nil, errors.Wrap(err, "inter.db.Find")
	}

	if len(events) == 0 {
		return nil, nil
	}

	return events, nil
}

// checkTimeIsFree проверяет, что указанное время свободно.
func (inter *Inter) checkTimeIsFree(ctx context.Context, event *domain.Event) error {
	filter := domain.Event{
		UserID:    event.UserID,
		StartDate: event.StartDate,
		EndDate:   event.EndDate,
	}

	events, err := inter.db.Find(ctx, domain.TakePeriodNotification, filter)
	if err != nil {
		return errors.Wrap(err, "inter.db.Find")
	}
	if len(events) > 0 {
		return domain.ErrDateBusy
	}

	return nil
}

// checkEventExists проверяет, что событие существует.
func (inter *Inter) checkEventExists(ctx context.Context, id string) error {
	filter := domain.Event{
		ID: id,
	}

	events, err := inter.db.Find(ctx, domain.TakePeriodNotification, filter)
	if err != nil {
		return errors.Wrap(err, "inter.db.Find")
	}
	if len(events) == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}
