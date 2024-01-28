package memorystorage

import (
	"context"
	"sync"

	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/common"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/domain"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/pkg/util"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/pkg/errors"
)

const (
	actionCreate = "CREATE"
	actionUpdate = "UPDATE"
	actionDelete = "DELETE"
	actionRead   = "READ"

	ctxEventID   = "event-id"
	ctxCondition = "condition-read"
	ctxStartDate = "start-date"
	ctxEndDate   = "end-date"
)

type Storage struct {
	data map[string]Event
	mu   sync.RWMutex
	log  common.Logger
}

func (s *Storage) Create(_ context.Context, in *domain.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := util.GenerateUUID()
	if err != nil {
		return "", errors.Wrap(err, "util.GenerateUUID")
	}
	in.ID = id

	event := eventFromDomain(in)
	dateEnd := &event.EndDate
	if event.EndDate.IsZero() {
		dateEnd = nil
	}
	s.log.Debugw(actionCreate, ctxEventID, id, ctxStartDate, event.StartDate, ctxEndDate, dateEnd)

	s.data[event.ID] = *event
	return event.ID, nil
}

func (s *Storage) Update(_ context.Context, in *domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dateEnd := &in.EndDate
	if in.EndDate.IsZero() {
		dateEnd = nil
	}
	s.log.Debugw(actionUpdate, ctxEventID, in.ID, ctxStartDate, in.StartDate, ctxEndDate, dateEnd)
	s.data[in.ID] = *eventFromDomain(in)
	return nil
}

func (s *Storage) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, id)
	s.log.Debugw(actionDelete, ctxEventID, id)
	return nil
}

func (s *Storage) Find(_ context.Context, condition int, filter domain.Event) ([]domain.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := make([]Event, 0, len(s.data))
	for _, value := range s.data {
		if (filter.StartDate.IsZero() && filter.EndDate.IsZero()) ||
			util.CheckIntersectDateIntervals(value.StartDate, value.EndDate, filter.StartDate, filter.EndDate) {
			events = append(events, value)
		}
	}

	s.log.Debugw(actionRead, ctxCondition, condition, ctxStartDate, filter.StartDate, ctxEndDate, filter)
	return eventsToDomain(events), nil
}

func New(log common.Logger) *Storage {
	return &Storage{
		data: map[string]Event{},
		mu:   sync.RWMutex{},
		log:  log,
	}
}
