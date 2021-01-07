//+build wireinject

package main

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	logrusadapter "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/logger/logrus"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/repository/memory"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/repository/sql"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/uuid"
	internalhttp "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/http/handler"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/usecase"
	"net/http"
)

func setup(cfg Config) (*CalendarApp, error) {
	wire.Build(
		NewApp,
		loggerProvider,

		httpServerProvider,
		httpHandlerProvider,
		handler.NewHandler,

		usecase.NewEventUseCase,
		eventRepositoryProvider,
		eventParticipantRepositoryProvider,
		dbProvider,
		sql.NewStatementBuilder,
		event.NewService,
		uuid.NewGenerator,
	)

	return nil, nil
}

func loggerProvider(cfg Config) (common.Logger, error) {
	return logrusadapter.New(logrusadapter.Config(cfg.Logger))
}

func httpServerProvider(cfg Config, handler http.Handler, l common.Logger) *internalhttp.Server {
	return internalhttp.NewServer(
		internalhttp.Config(cfg.Server),
		handler,
		l,
	)
}

func httpHandlerProvider(h *handler.Handler) http.Handler {
	return h.InitRoutes()
}

func eventRepositoryProvider(
	cfg Config,
	participantRepository event.ParticipantRepository,
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

func eventParticipantRepositoryProvider(cfg Config, db *sqlx.DB, sbt sq.StatementBuilderType) (event.ParticipantRepository, error) {
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
