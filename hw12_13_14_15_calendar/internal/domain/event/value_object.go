package event

import (
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
	"time"
)

type Interval struct {
	BeginDate time.Time
	EndDate   time.Time
}

type UserInterval struct {
	Uid user.UID
	Interval
}
