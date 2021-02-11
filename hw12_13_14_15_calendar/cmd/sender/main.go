package main

import (
	"context"
	"flag"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/application/broker"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/spf13/viper"
	golog "log"
	"os"
	"os/signal"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/sender/config.toml", "Path to configuration file")
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
	flag.Parse()

	// parsing config
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		golog.Fatal(err)
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		golog.Fatal(err)
	}

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

		notificationCh, err := app.Reader.Read(ctx)
		if err != nil {
			app.Logger.Error(err.Error())
			cancel()
			return
		}

		for notification := range notificationCh {
			str := string(notification)
			golog.Println(str)
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
