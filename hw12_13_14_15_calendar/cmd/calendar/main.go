package main

import (
	"context"
	"errors"
	"flag"
	golog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	internalgrpc "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

type CalendarApp struct {
	logger     common.Logger
	httpServer *internalhttp.Server
	grpcServer *internalgrpc.Server
}

func NewApp(logger common.Logger, httpServer *internalhttp.Server, grpcServer *internalgrpc.Server) *CalendarApp {
	return &CalendarApp{logger: logger, httpServer: httpServer, grpcServer: grpcServer}
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
	go signalHandler(ctx, app, cancel)

	// start
	app.logger.Info("calendar is running...")

	go func() {
		if err := app.httpServer.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				app.logger.Info("http server stopped")
			} else {
				app.logger.Error("failed to start http server: " + err.Error())
				cancel()
			}
		}
	}()

	go func() {
		if err := app.grpcServer.Start(); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				app.logger.Info("grpc server stopped")
			} else {
				app.logger.Error("failed to start grpc server: " + err.Error())
			}
		}
	}()

	<-ctx.Done()
}

func signalHandler(ctx context.Context, app *CalendarApp, cancel context.CancelFunc) {
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	select {
	case <-signals:
		signal.Stop(signals)

		serverCloseCtx, serverCancel := context.WithTimeout(context.Background(), time.Second*3)
		defer serverCancel()

		if err := app.httpServer.Stop(serverCloseCtx); err != nil {
			app.logger.Error("failed to stop http server: " + err.Error())
		}

		app.grpcServer.Stop()

	case <-ctx.Done():
	}
}
