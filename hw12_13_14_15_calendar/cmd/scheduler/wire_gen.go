// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/application/scheduler"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/logger/logrus"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/rabbit"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/repository/memory"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/repository/sql"
	"time"
)

// Injectors from wire.go:

func setup(cfg Config) (*SchedulerApp, func(), error) {
	logger, err := loggerProvider(cfg)
	if err != nil {
		return nil, nil, err
	}
	eventParticipantRepository := memory.NewEventParticipantRepository()
	db, err := dbProvider(cfg)
	if err != nil {
		return nil, nil, err
	}
	statementBuilderType := sql.NewStatementBuilder()
	repository, err := eventRepositoryProvider(cfg, eventParticipantRepository, db, statementBuilderType)
	if err != nil {
		return nil, nil, err
	}
	cleanerImpl, err := cleanerProvider(cfg, repository)
	if err != nil {
		return nil, nil, err
	}
	client := rabbitClientProvider(cfg)
	participantRepository := sql.NewEventParticipantRepository(db, statementBuilderType)
	notifierImpl := scheduler.NewNotifierImpl(client, repository, participantRepository)
	schedulerApp := NewApp(logger, cleanerImpl, notifierImpl)
	return schedulerApp, func() {
	}, nil
}

// wire.go:

func loggerProvider(cfg Config) (common.Logger, error) {
	return logrusadapter.New(logrusadapter.Config(cfg.Logger))
}

func eventRepositoryProvider(
	cfg Config,
	participantRepository *memory.EventParticipantRepository,
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

func dbProvider(cfg Config) (*sqlx.DB, error) {
	if cfg.Repository.Type != "sql" {
		return nil, nil
	}

	return sql.NewDB(cfg.DB.DSN)
}

func cleanerProvider(cfg Config, repository event.Repository) (*scheduler.CleanerImpl, error) {
	lifespan, err := time.ParseDuration(cfg.Cleaner.EventLifespan)
	if err != nil {
		return nil, err
	}

	return scheduler.NewCleanerImpl(lifespan, repository), nil
}

func rabbitClientProvider(cfg Config) *rabbit.Client {
	return rabbit.NewClient(rabbit.Config(cfg.Rabbit))
}