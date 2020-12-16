package main

import (
	"context"
	"errors"
	"flag"
	golog "log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/logger"
	logrusadapter "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/logger/logrus"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/repository/memory"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/repository/sql"
	internalhttp "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/server/http/handler"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/usecase"
	"github.com/spf13/viper"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
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

	// build dependencies
	l, err := logrusadapter.New(logrusadapter.Config(cfg.Logger))
	if err != nil {
		golog.Fatal(err)
	}

	repoFactory, err := repository.GetFactory(
		cfg.Repository.Type,
		new(memory.Factory),
		new(sql.Factory),
	)
	if err != nil {
		golog.Fatal(err)
	}

	repos, err := repoFactory.Build(cfg.Repository.Credentials)
	if err != nil {
		golog.Fatal(err)
	}

	useCases := usecase.NewUseCase(repos)
	h := handler.NewHandler(useCases, l)
	server := internalhttp.NewServer(cfg.Server.Port, h.InitRoutes(), l)

	// handle os signals
	ctx, cancel := context.WithCancel(context.Background())
	go signalHandler(server, l, cancel)

	// start
	l.Info("calendar is running...")
	if err := server.Start(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			l.Info("calendar stopped")
		} else {
			l.Error("failed to start http server: " + err.Error())
			os.Exit(1)
		}
	}

	<-ctx.Done()
}

func signalHandler(s *internalhttp.Server, l logger.Logger, cancel context.CancelFunc) {
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals)

	<-signals
	signal.Stop(signals)

	ctx, serverCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer serverCancel()

	if err := s.Stop(ctx); err != nil {
		l.Error("failed to stop http server: " + err.Error())
	}
}
