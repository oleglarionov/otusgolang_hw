package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/domain/user"
	"net/http"
)

func authMiddleware(c *gin.Context) {
	value := c.Request.Header.Get("x-uid")
	if value == "" {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	c.Set("uid", user.UID(value))
	c.Next()
}
