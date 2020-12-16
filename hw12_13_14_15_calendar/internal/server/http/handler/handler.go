package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/usecase"
)

type Handler struct {
	useCases *usecase.UseCase
	l        logger.Logger
}

func NewHandler(useCases *usecase.UseCase, l logger.Logger) *Handler {
	return &Handler{
		useCases: useCases,
		l:        l,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()

	r.GET("/ping", func(ctx *gin.Context) {
		_, err := ctx.Writer.Write([]byte("pong\n"))
		if err != nil {
			h.l.Error("error writing response: " + err.Error())
		}
	})

	return r
}
