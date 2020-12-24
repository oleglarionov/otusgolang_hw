package usecase

import "github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/repository"

type UseCase struct {
	repos *repository.Repository
}

func NewUseCase(repos *repository.Repository) *UseCase {
	return &UseCase{repos: repos}
}
