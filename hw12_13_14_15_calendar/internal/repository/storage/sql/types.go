package dbstorage

import (
	"database/sql"
	"time"
)

// Event событие в календаре.
type Event struct {
	ID           string        `db:"id"`           // Идентификатор события;
	UserID       int64         `db:"user_id"`      // Идентификатор пользователя, владельца события
	Title        string        `db:"title"`        // Заголовок
	StartDate    time.Time     `db:"start_date"`   // Дата и время начала события
	EndDate      time.Time     `db:"end_date"`     // Дата и время окончания
	Notification time.Duration `db:"notification"` // Оповещение о событии, за какое время необходимо отправить
	// уведомление пользователю
	Description sql.NullString `db:"description"` // Описание события
}
