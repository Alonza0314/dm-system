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

// withLogging wraps a route's handler so the request is logged exactly once,
// after it completes, with the level derived from the actual response status
// code the handler wrote (via c.JSON/c.Status). This replaces the manual
// "<name> request from ..." / "<name> successful from ..." pair that used to
// open and close every handler -- register it once per route instead:
//
//	HandlerFunc: withLogging("GetCategories", b.CatLog, b.handleGetCategories),
func withLogging(name string, lg loggergoModel.LoggerInterface, handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c)

		status := c.Writer.Status()
		switch {
		case status >= http.StatusInternalServerError:
			lg.Errorf("%s failed (status %d) for %s", name, status, c.ClientIP())
		case status >= http.StatusBadRequest:
			lg.Warnf("%s failed (status %d) for %s", name, status, c.ClientIP())
		default:
			lg.Infof("%s succeeded (status %d) for %s", name, status, c.ClientIP())
		}
	}
}
