package main

import (
	"context"
	golog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/application/broker"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
}

type SenderApp struct {
	Logger common.Logger
	Reader broker.Reader
}

func NewApp(logger common.Logger, reader broker.Reader) *SenderApp {
	return &SenderApp{
		Logger: logger,
		Reader: reader,
	}
}

func main() {
	cfg := getConfig()

	// setup app
	app, cleanup, err := setup(cfg)
	if err != nil {
		golog.Fatal(err)
	}
	defer cleanup()

	// handle os signals
	ctx, cancel := context.WithCancel(context.Background())
	go signalHandler(ctx, cancel)

	go func() {
		app.Logger.Info("sender started")
		defer app.Logger.Info("sender stopped")
		defer cancel()

		notificationCh, err := app.Reader.Read(ctx)
		if err != nil {
			app.Logger.Error(err.Error())
			return
		}

		for notification := range notificationCh {
			str := string(notification)
			golog.Println(str)
			app.Logger.Info(str)
		}
	}()

	<-ctx.Done()
}

func signalHandler(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	defer signal.Stop(signals)

	select {
	case <-signals:
	case <-ctx.Done():
	}
}

func getConfig() Config {
	return Config{
		Logger: LoggerConf{
			Level: viper.GetString("LOG_LEVEL"),
			File:  viper.GetString("LOG_FILE"),
		},
		Rabbit: RabbitConf{
			DSN:   viper.GetString("RABBIT_DSN"),
			Queue: viper.GetString("RABBIT_QUEUE"),
		},
	}
}
