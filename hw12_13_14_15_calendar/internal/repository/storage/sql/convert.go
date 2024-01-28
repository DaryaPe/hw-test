package dbstorage

import (
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/domain"
)

func eventFromDomain(in *domain.Event) *Event {
	return &Event{
		ID:           in.ID,
		UserID:       in.UserID,
		Title:        in.Title,
		StartDate:    in.StartDate,
		EndDate:      in.EndDate,
		Notification: in.Notification,
		Description:  StringToNull(in.Description),
	}
}

func eventToDomain(in *Event) *domain.Event {
	return &domain.Event{
		ID:           in.ID,
		UserID:       in.UserID,
		Title:        in.Title,
		StartDate:    in.StartDate,
		EndDate:      in.EndDate,
		Notification: in.Notification,
		Description:  in.Description.String,
	}
}

func eventsFromDomain(in []Event) []domain.Event {
	list := make([]domain.Event, 0, len(in))
	for i := range in {
		list = append(list, *eventToDomain(&in[i]))
	}
	return list
}
