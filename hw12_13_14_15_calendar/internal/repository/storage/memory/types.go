package memorystorage

import (
	"time"
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
