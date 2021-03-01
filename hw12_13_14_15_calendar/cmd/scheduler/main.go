package main

import (
	"context"
	golog "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/application/scheduler"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
}

type SchedulerApp struct {
	Logger   common.Logger
	Cleaner  scheduler.Cleaner
	Notifier scheduler.Notifier
}

func NewApp(logger common.Logger, cleaner scheduler.Cleaner, notifier scheduler.Notifier) *SchedulerApp {
	return &SchedulerApp{Logger: logger, Cleaner: cleaner, Notifier: notifier}
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
		cleanerTicker := time.NewTicker(5 * time.Minute)
		notifierTicker := time.NewTicker(3 * time.Second)
		app.Logger.Info("scheduler started")
		defer app.Logger.Info("scheduler stopped")
		for {
			select {
			case <-cleanerTicker.C:
				err = app.Cleaner.Clean(ctx)
				if err != nil {
					app.Logger.Error(err.Error())
				}
			case <-notifierTicker.C:
				err = app.Notifier.Notify(ctx)
				if err != nil {
					app.Logger.Error(err.Error())
				}
			case <-ctx.Done():
				return
			}
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
			Level: viper.Get("LOG_LEVEL").(string),
			File:  viper.Get("LOG_FILE").(string),
		},
		DB: DBConf{
			DSN: viper.Get("DB_DSN").(string),
		},
		Repository: RepositoryConf{
			Type: viper.Get("REPO_TYPE").(string),
		},
		Cleaner: CleanerConf{
			EventLifespan: viper.Get("EVENT_LIFESPAN").(string),
		},
		Rabbit: RabbitConf{
			DSN:   viper.Get("RABBIT_DSN").(string),
			Queue: viper.Get("RABBIT_QUEUE").(string),
		},
	}
}
