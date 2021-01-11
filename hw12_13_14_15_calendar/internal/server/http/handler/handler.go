package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/usecase"
)

type Handler struct {
	eventUseCase usecase.EventUseCase
	l            common.Logger
}

func NewHandler(
	eventUseCase usecase.EventUseCase,
	l common.Logger,
) *Handler {
	return &Handler{
		eventUseCase: eventUseCase,
		l:            l,
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

	api := r.Group("/api")
	api.Use(authMiddleware)
	{
		events := api.Group("/events")
		{
			events.POST("/", h.createEvent)
			events.PUT("/:id", h.updateEvent)
			events.DELETE("/:id", h.deleteEvent)
			events.GET("/day/:date", h.dayEvents)
			events.GET("/week/:beginDate", h.weekEvents)
			events.GET("/month/:beginDate", h.monthEvents)
		}
	}

	return r
}
