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
			HandlerFunc: b.handleBorrow,
		},
		{
			Name:        "Return",
			Method:      http.MethodDelete,
			Pattern:     "/qrcode/:cate/:dev",
			HandlerFunc: b.handleReturn,
		},
	}
}

func (b *backend) handleBorrow(c *gin.Context) {
	b.QrdLog.Infof("Borrow %s %s fromo %s", c.Param("cate"), c.Param("dev"), c.ClientIP())

	var req model.RequestQrcodeBorrow
	if !requestBinding(c, &req) {
		b.QrdLog.Warnf("Invalid qrcode borrow request from %s", c.ClientIP())
		return
	}

	if errDetail := b.Processor.Borrow(c.Param("cate"), c.Param("dev"), &req); errDetail != nil {
		errDetailLog(errDetail, b.QrdLog, "Qrcode borrow failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	b.QrdLog.Infof("Borrow %s %s successful from %s", c.Param("cate"), c.Param("dev"), c.ClientIP())
	c.Status(http.StatusOK)
}

func (b *backend) handleReturn(c *gin.Context) {
	b.QrdLog.Infof("Return %s %s fromo %s", c.Param("cate"), c.Param("dev"), c.ClientIP())

	var req model.RequestQrcodeReturn
	if !requestBinding(c, &req) {
		b.QrdLog.Warnf("Invalid qrcode return request from %s", c.ClientIP())
		return
	}

	if errDetail := b.Processor.Return(c.Param("cate"), c.Param("dev"), &req); errDetail != nil {
		errDetailLog(errDetail, b.QrdLog, "Qrcode return failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	b.QrdLog.Infof("Return %s %s successful from %s", c.Param("cate"), c.Param("dev"), c.ClientIP())
	c.Status(http.StatusOK)
}
