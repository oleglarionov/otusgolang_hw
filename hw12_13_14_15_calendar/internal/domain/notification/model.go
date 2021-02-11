package notification

import (
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
	"time"
)

type Model struct {
	Uid   user.UID
	Event Event
}

type Event struct {
	Title       string
	Description string
	BeginDate   time.Time
	EndDate     time.Time
}
