package middleware

import (
	"github.com/gin-gonic/gin"
	"go-sever/common"
	"go-sever/component"
	"go-sever/library/log"
	"net/http"
)

func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			component.ErrLogger.Error(log.F{
				"log_type": common.LogTypeForPanic,
				"info":     err,
			})
			c.Status(http.StatusInternalServerError)
		}
	}()

	c.Next()
}
