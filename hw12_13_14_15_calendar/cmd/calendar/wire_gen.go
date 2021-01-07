// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/logger/logrus"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/repository/memory"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/repository/sql"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/uuid"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/http/handler"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/usecase"
	"net/http"
)

// Injectors from wire.go:

func setup(cfg Config) (*CalendarApp, error) {
	logger, err := loggerProvider(cfg)
	if err != nil {
		return nil, err
	}
	db, err := dbProvider(cfg)
	if err != nil {
		return nil, err
	}
	statementBuilderType := sql.NewStatementBuilder()
	participantRepository, err := eventParticipantRepositoryProvider(cfg, db, statementBuilderType)
	if err != nil {
		return nil, err
	}
	repository, err := eventRepositoryProvider(cfg, participantRepository, db, statementBuilderType)
	if err != nil {
		return nil, err
	}
	service := event.NewService(repository, participantRepository)
	uuidGenerator := uuid.NewGenerator()
	eventUseCaseInterface := usecase.NewEventUseCase(repository, participantRepository, service, uuidGenerator)
	handlerHandler := handler.NewHandler(eventUseCaseInterface, logger)
	httpHandler := httpHandlerProvider(handlerHandler)
	server := httpServerProvider(cfg, httpHandler, logger)
	calendarApp := NewApp(logger, server)
	return calendarApp, nil
}

// wire.go:

func loggerProvider(cfg Config) (common.Logger, error) {
	return logrusadapter.New(logrusadapter.Config(cfg.Logger))
}

func httpServerProvider(cfg Config, handler2 http.Handler, l common.Logger) *internalhttp.Server {
	return internalhttp.NewServer(internalhttp.Config(cfg.Server), handler2, l)
}

func httpHandlerProvider(h *handler.Handler) http.Handler {
	return h.InitRoutes()
}

func eventRepositoryProvider(
	cfg Config,
	participantRepository event.ParticipantRepository,
	db *sqlx.DB,
	sbt squirrel.StatementBuilderType,
) (event.Repository, error) {
	repoType := cfg.Repository.Type
	switch repoType {
	case "memory":
		return memory.NewEventRepository(participantRepository), nil
	case "sql":
		return sql.NewEventRepository(db, sbt), nil
	default:
		return nil, fmt.Errorf("unsupported participantRepository type: %s", repoType)
	}
}

func eventParticipantRepositoryProvider(cfg Config, db *sqlx.DB, sbt squirrel.StatementBuilderType) (event.ParticipantRepository, error) {
	repoType := cfg.Repository.Type
	switch repoType {
	case "memory":
		return memory.NewEventParticipantRepository(), nil
	case "sql":
		return sql.NewEventParticipantRepository(db, sbt), nil
	default:
		return nil, fmt.Errorf("unsupported repository type: %s", repoType)
	}
}

func dbProvider(cfg Config) (*sqlx.DB, error) {
	if cfg.Repository.Type != "sql" {
		return nil, nil
	}

	return sql.NewDb(cfg.DB.DSN)
}