//+build wireinject

package main

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/application/broker"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/application/scheduler"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	logrusadapter "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/logger/logrus"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/rabbit"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/repository/memory"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/repository/sql"
)

func setup(cfg Config) (*SchedulerApp, func(), error) {
	wire.Build(
		NewApp,
		loggerProvider,
		wire.Bind(new(scheduler.Cleaner), new(*scheduler.CleanerImpl)),
		cleanerProvider,
		wire.Bind(new(scheduler.Notifier), new(*scheduler.NotifierImpl)),
		scheduler.NewNotifierImpl,
		wire.Bind(new(broker.Pusher), new(*rabbit.Client)),
		rabbitClientProvider,
		eventRepositoryProvider,
		memory.NewEventParticipantRepository,
		sql.NewEventParticipantRepository,
		dbProvider,
		sql.NewStatementBuilder,
	)

	return nil, nil, nil
}

func loggerProvider(cfg Config) (common.Logger, error) {
	return logrusadapter.New(logrusadapter.Config(cfg.Logger))
}

func eventRepositoryProvider(
	cfg Config,
	participantRepository *memory.EventParticipantRepository,
	db *sqlx.DB,
	sbt sq.StatementBuilderType,
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

func rabbitClientProvider(cfg Config) (*rabbit.Client, func()) {
	client := rabbit.NewClient(rabbit.Config(cfg.Rabbit))
	cleanup := func() {
		client.Close()
	}
	return client, cleanup
}
