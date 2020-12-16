package sql

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	// init driver.
	_ "github.com/lib/pq"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/repository"
)

type Factory struct {
}

type Config struct {
	Dsn string `mapstructure:"dsn"`
}

var _ repository.RepoFactory = (*Factory)(nil)

func (f *Factory) Build(dsn interface{}) (*repository.Repository, error) {
	db, err := sqlx.Connect("postgres", dsn.(string))
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &repository.Repository{
		EventRepository: NewEventRepository(db),
	}, nil
}

func (f *Factory) RepoType() repository.RepoType {
	return "sql"
}
