package internal

import (
	"backend/model"
	"net/http"

	loggergoModel "github.com/Alonza0314/logger-go/v2/model"
	"github.com/gin-gonic/gin"
)

func requestBinding(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return false
	}
	return true
}

func errDetailLog(errorDetail *model.ErrorDetail, lg loggergoModel.LoggerInterface, format string, args ...interface{}) {
	switch errorDetail.HttpStatus {
	case http.StatusInternalServerError:
		lg.Errorf(format, args...)
	default:
		lg.Warnf(format, args...)
	}
}
