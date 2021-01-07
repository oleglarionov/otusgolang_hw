package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/usecase"
	"net/http"
	"time"
)

func (h *Handler) createEvent(c *gin.Context) {
	uid := c.MustGet("uid").(user.UID)

	var input usecase.CreateEventDto
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.eventUseCase.Create(c, uid, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) updateEvent(c *gin.Context) {
	uid := c.MustGet("uid").(user.UID)
	id := event.ID(c.Param("id"))

	var input usecase.UpdateEventDto
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.eventUseCase.Update(c, uid, id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) deleteEvent(c *gin.Context) {
	uid := c.MustGet("uid").(user.UID)
	id := event.ID(c.Param("id"))

	if err := h.eventUseCase.Delete(c, uid, id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}

func (h *Handler) dayEvents(c *gin.Context) {
	uid := c.MustGet("uid").(user.UID)

	strDate := c.Param("date")
	day, err := time.Parse("2006-01-02", strDate)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.eventUseCase.DayList(c, uid, day)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) weekEvents(c *gin.Context) {
	uid := c.MustGet("uid").(user.UID)

	strDate := c.Param("beginDate")
	beginDate, err := time.Parse("2006-01-02", strDate)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.eventUseCase.WeekList(c, uid, beginDate)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) monthEvents(c *gin.Context) {
	uid := c.MustGet("uid").(user.UID)

	strDate := c.Param("beginDate")
	beginDate, err := time.Parse("2006-01-02", strDate)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.eventUseCase.MonthList(c, uid, beginDate)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, result)
}
