package internal

import (
	"backend/model"
	"net/http"

	"github.com/free-ran-ue/util"
	"github.com/gin-gonic/gin"
)

func (b *backend) getQrcodeRoutes() util.Routes {
	return util.Routes{
		{
			Name:        "Borrow",
			Method:      http.MethodPost,
			Pattern:     "/qrcode/:cate/:dev",
			HandlerFunc: withLogging("Borrow", b.QrdLog, b.handleBorrow),
		},
		{
			Name:        "Return",
			Method:      http.MethodDelete,
			Pattern:     "/qrcode/:cate/:dev",
			HandlerFunc: withLogging("Return", b.QrdLog, b.handleReturn),
		},
	}
}

func (b *backend) handleBorrow(c *gin.Context) {
	var req model.RequestQrcodeBorrow
	if !requestBinding(c, &req) {
		return
	}

	if errDetail := b.Processor.Borrow(c.Param("cate"), c.Param("dev"), &req); errDetail != nil {
		errDetailLog(errDetail, b.QrdLog, "Qrcode borrow failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.Status(http.StatusOK)
}

func (b *backend) handleReturn(c *gin.Context) {
	var req model.RequestQrcodeReturn
	if !requestBinding(c, &req) {
		return
	}

	if errDetail := b.Processor.Return(c.Param("cate"), c.Param("dev"), &req); errDetail != nil {
		errDetailLog(errDetail, b.QrdLog, "Qrcode return failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.Status(http.StatusOK)
}
