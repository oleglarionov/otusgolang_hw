package memory

import "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/repository"

type Factory struct {
}

var _ repository.RepoFactory = (*Factory)(nil)

func (f *Factory) Build(_ interface{}) (*repository.Repository, error) {
	return &repository.Repository{
		EventRepository: NewEventRepo(),
	}, nil
}

func (f *Factory) RepoType() repository.RepoType {
	return "memory"
}
