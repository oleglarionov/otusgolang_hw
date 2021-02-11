//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/application/broker"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	logrusadapter "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/logger/logrus"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/infrstructure/rabbit"
)

func setup(cfg Config) (*SenderApp, func(), error) {
	wire.Build(
		NewApp,
		loggerProvider,
		wire.Bind(new(broker.Reader), new(*rabbit.Client)),
		rabbitClientProvider,
	)

	return nil, nil, nil
}

func loggerProvider(cfg Config) (common.Logger, error) {
	return logrusadapter.New(logrusadapter.Config(cfg.Logger))
}

func rabbitClientProvider(cfg Config) (*rabbit.Client, func()) {
	client := rabbit.NewClient(rabbit.Config(cfg.Rabbit))
	cleanup := func() {
		client.Close()
	}
	return client, cleanup
}
