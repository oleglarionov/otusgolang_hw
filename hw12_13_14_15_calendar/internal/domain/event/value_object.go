package event

import (
	"time"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
)

type Interval struct {
	BeginDate time.Time
	EndDate   time.Time
}

type UserInterval struct {
	UID user.UID
	Interval
}
