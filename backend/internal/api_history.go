package internal

import (
	"net/http"

	"github.com/free-ran-ue/util"
	"github.com/gin-gonic/gin"
)

func (b *backend) getHistoryRoutes() util.Routes {
	return util.Routes{
		{
			Name:        "GetHistory",
			Method:      http.MethodGet,
			Pattern:     "/history/:cate/:dev",
			HandlerFunc: withLogging("GetHistory", b.HstLog, b.handleGetHistory),
		},
	}
}

func (b *backend) handleGetHistory(c *gin.Context) {
	response, errDetail := b.Processor.GetHistory(c.Param("cate"), c.Param("dev"))
	if errDetail != nil {
		errDetailLog(errDetail, b.HstLog, "Get history failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.JSON(http.StatusOK, response)
}
