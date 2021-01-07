package main

import (
	"context"
	"errors"
	"flag"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"net/http"

	golog "log"
	"os"
	"os/signal"
	"time"

	internalhttp "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/spf13/viper"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

type CalendarApp struct {
	logger     common.Logger
	httpServer *internalhttp.Server
}

func NewApp(logger common.Logger, httpServer *internalhttp.Server) *CalendarApp {
	return &CalendarApp{logger: logger, httpServer: httpServer}
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

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
	app, err := setup(cfg)
	if err != nil {
		golog.Fatal(err)
	}

	// handle os signals
	ctx, cancel := context.WithCancel(context.Background())
	go signalHandler(app, ctx, cancel)

	// start
	app.logger.Info("calendar is running...")
	if err := app.httpServer.Start(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			app.logger.Info("calendar stopped")
		} else {
			app.logger.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}

	<-ctx.Done()
}

func signalHandler(app *CalendarApp, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals)

	select {
	case <-signals:
		signal.Stop(signals)

		serverCloseCtx, serverCancel := context.WithTimeout(context.Background(), time.Second*3)
		defer serverCancel()

		if err := app.httpServer.Stop(serverCloseCtx); err != nil {
			app.logger.Error("failed to stop http server: " + err.Error())
		}

	case <-ctx.Done():
	}
}
