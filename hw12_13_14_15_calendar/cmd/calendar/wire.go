//+build wireinject

package main

import (
	"fmt"
	"net/http"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/api"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	logrusadapter "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/logger/logrus"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/repository/memory"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/repository/sql"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/uuid"
	internalgrpc "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/grpc"
	grpchandler "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/grpc/handler"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/grpc/middleware"
	internalhttp "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/http/handler"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/usecase"
)

func setup(cfg Config) (*CalendarApp, error) {
	wire.Build(
		NewApp,
		loggerProvider,

		httpServerProvider,
		httpHandlerProvider,
		handler.NewHandler,

		grpcServerProvider,
		wire.Bind(new(api.EventServiceServer), new(*grpchandler.EventServiceServerImpl)),
		grpchandler.NewEventServiceServerImpl,
		grpcMiddlewaresProvider,
		middleware.NewLoggingMiddleware,
		middleware.NewAuthenticationMiddleware,

		wire.Bind(new(usecase.EventUseCase), new(*usecase.EventUseCaseImpl)),
		usecase.NewEventUseCaseImpl,
		locationProvider,
		eventRepositoryProvider,

		memory.NewEventParticipantRepository,
		sql.NewEventParticipantRepository,

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
		internalhttp.Config(cfg.HTTPServer),
		handler,
		l,
	)
}

func httpHandlerProvider(h *handler.Handler) http.Handler {
	return h.InitRoutes()
}

func grpcServerProvider(cfg Config, eventService api.EventServiceServer, middlewares []middleware.Middleware) (*internalgrpc.Server, error) {
	return internalgrpc.NewServer(cfg.GrpcServer.Port, eventService, middlewares)
}

func grpcMiddlewaresProvider(
	loggingMiddleware *middleware.LoggingMiddleware,
	authenticationMiddleware *middleware.AuthenticationMiddleware,
) []middleware.Middleware {
	return []middleware.Middleware{
		loggingMiddleware,
		authenticationMiddleware,
	}
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

func locationProvider() (*time.Location, error) {
	return time.LoadLocation("Europe/Moscow")
}
