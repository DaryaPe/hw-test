package dbstorage

import (
	"context"
	"database/sql"
	"strings"

	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/domain"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/internal/pkg/util"
	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/pkg/errors"
)

type Row interface {
	Scan(dest ...interface{}) (err error)
	Err() error
}

type Rows interface {
	Scan(dest ...interface{}) (err error)
	Next() bool
	Err() error
	Close() error
}

type DB interface {
	Open(ctx context.Context) error
	Close() error
	Query(ctx context.Context, sql string, args ...interface{}) (context.CancelFunc, *sql.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) (context.CancelFunc, *sql.Row)
	Exec(ctx context.Context, sql string, args ...interface{}) error
	Rebind(query string) string
}

type Storage struct {
	db DB
}

func New(db DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Create(ctx context.Context, event *domain.Event) (string, error) {
	id, err := util.GenerateUUID()
	if err != nil {
		return "", errors.Wrap(err, "util.GenerateUUID")
	}
	event.ID = id

	dbEvent := eventFromDomain(event)
	query := `INSERT INTO 
			calendar.event (id, user_id, title, start_date, end_date, notification, description) 
			VALUES (:id, :user_id, :title, :start_date, :end_date, :notification, :description)`

	err = s.db.Exec(ctx, query, dbEvent)
	if err != nil {
		return "", errors.Wrap(err, "s.db.Exec")
	}

	return event.ID, nil
}

func (s *Storage) Update(ctx context.Context, event *domain.Event) error {
	dEvent := eventFromDomain(event)
	query := `UPDATE calendar.event 
	SET user_id = :user_id, 
	    title = :title, 
	    start_date = :start_date, 
	    end_date = :end_date, 
	    notification = :notification, 
	    description = description
	WHERE  id = :id`

	err := s.db.Exec(ctx, query, dEvent)
	if err != nil {
		return errors.Wrap(err, "s.db.Exec")
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM calendar.event WHERE id = $1`

	if err := s.db.Exec(ctx, query, id); err != nil {
		return errors.Wrap(err, "s.db.Exec")
	}
	return nil
}

func (s *Storage) Find(ctx context.Context, _ int, filter domain.Event) ([]domain.Event, error) {
	var query strings.Builder
	query.WriteString(
		`SELECT id, user_id, title, start_date, end_date, notification, description FROM calendar.event where 1=1`)

	var args []interface{}

	if !filter.StartDate.IsZero() {
		query.WriteString(` and start_date >= ?`)
		args = append(args, filter.StartDate)
	}
	if !filter.EndDate.IsZero() {
		query.WriteString(` AND end_date <= ?`)
		args = append(args, filter.EndDate)
	}
	if len(filter.ID) > 0 {
		query.WriteString(` AND id = ?`)
		args = append(args, filter.ID)
	}
	if filter.UserID > 0 {
		query.WriteString(` AND user_id = ?`)
		args = append(args, filter.UserID)
	}

	cancel, rows, err := s.db.Query(ctx, s.db.Rebind(query.String()), args...)
	if err != nil {
		return nil, errors.Wrap(err, "s.db.Query")
	}
	defer cancel()
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}

	var events []Event
	for rows.Next() {
		var event Event
		if err = rows.Scan(
			&event.ID,
			&event.UserID,
			&event.Title,
			&event.StartDate,
			&event.EndDate,
			&event.Notification,
			&event.Description,
		); err != nil {
			return nil, errors.Wrap(err, "rows.Next")
		}
		events = append(events, event)
	}

	return eventsFromDomain(events), nil
}

func (s *Storage) Connect(ctx context.Context) error {
	if err := s.db.Open(ctx); err != nil {
		return errors.Wrap(err, "s.db.Open")
	}
	return nil
}

func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return errors.Wrap(err, "s.db.Close")
	}
	return nil
}
