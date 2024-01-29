package domain

import (
	"context"
	"errors"
	"time"
)

const (
	AllNotifications = iota
	TakePeriodNotification
	TakeDayPeriodNotification
	TakeWeekPeriodNotification
	TakeMonthPeriodNotification
)

var (
	ErrNotDefinedPeriod = errors.New("not defined period for getting notifications")
	ErrDateBusy         = errors.New("selected date already busy")
	ErrEventNotFound    = errors.New("event not found")
)

// Event событие в календаре.
type Event struct {
	ID           string        // Идентификатор события;
	UserID       int64         // Идентификатор пользователя, владельца события
	Title        string        // Заголовок
	StartDate    time.Time     // Дата и время начала события
	EndDate      time.Time     // Дата и время окончания
	Notification time.Duration // Оповещение о событии, за какое время необходимо отправить уведомление пользователю
	Description  string        // Описание события
}

// GetNotification возвращает уведомление на основании события.
func (e Event) GetNotification() Notification {
	return Notification{
		ID:        e.ID,
		UserID:    e.UserID,
		Title:     e.Title,
		StartDate: e.StartDate,
	}
}

// Notification оповещение пользователя о событии.
type Notification struct {
	ID        string    // Идентификатор события
	UserID    int64     // Пользователь, которому необходимо отправить уведомление
	Title     string    // Заголовок события
	StartDate time.Time // Дата начала события
}

type EventInter interface {
	Create(ctx context.Context, event *Event) (string, error)
	Update(ctx context.Context, event *Event) error
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, date time.Time, condition int) ([]Event, error)
}
