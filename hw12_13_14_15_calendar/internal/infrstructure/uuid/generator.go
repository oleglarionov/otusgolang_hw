package uuid

import (
	"github.com/google/uuid"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
)

func NewGenerator() common.UUIDGenerator {
	return &generator{}
}

type generator struct {
}

func (g *generator) Generate() string {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	return newUUID.String()
}
