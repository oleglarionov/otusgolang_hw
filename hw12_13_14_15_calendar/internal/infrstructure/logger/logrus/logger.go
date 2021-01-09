package logrusadapter

import (
	"fmt"
	"os"
	"strings"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/sirupsen/logrus"
)

type Adapter struct {
	l *logrus.Logger
}

type Config struct {
	Level string
	File  string
}

var _ common.Logger = (*Adapter)(nil)

func New(cfg Config) (common.Logger, error) {
	l := logrus.New()

	file, err := os.OpenFile(cfg.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", cfg.File, err)
	}
	l.SetOutput(file)

	level, ok := levelMap[strings.ToLower(cfg.Level)]
	if !ok {
		return nil, fmt.Errorf("unknown level: %s", cfg.Level)
	}
	l.SetLevel(level)

	return &Adapter{l: l}, nil
}

func (a *Adapter) Error(msg string) {
	a.l.Error(msg)
}

func (a *Adapter) Warn(msg string) {
	a.l.Warn(msg)
}

func (a *Adapter) Info(msg string) {
	a.l.Info(msg)
}

func (a *Adapter) Debug(msg string) {
	a.l.Debug(msg)
}

var levelMap = map[string]logrus.Level{
	"error": logrus.ErrorLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
}
