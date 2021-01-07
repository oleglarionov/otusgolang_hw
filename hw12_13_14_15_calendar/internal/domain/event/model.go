package event

import (
	"database/sql/driver"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
	"time"
)

type Model struct {
	Id          ID        `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	BeginDate   time.Time `db:"begin_date"`
	EndDate     time.Time `db:"end_date"`
}

type Participant struct {
	EventId ID       `db:"event_id"`
	Uid     user.UID `db:"uid"`
}

type ID string

func (id ID) Value() (driver.Value, error) {
	return string(id), nil
}

func (id *ID) Scan(src interface{}) error {
	*id = ID(src.([]byte))
	return nil
}
