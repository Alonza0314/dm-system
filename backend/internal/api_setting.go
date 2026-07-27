package internal

import (
	"backend/model"
	"net/http"

	"github.com/free-ran-ue/util"
	"github.com/gin-gonic/gin"
)

func (b *backend) getSettingRoutes() util.Routes {
	return util.Routes{
		{
			Name:        "Setting account",
			Method:      http.MethodPost,
			Pattern:     "/setting/account",
			HandlerFunc: withLogging("Setting account", b.SetLog, b.handleSettingAccount),
		},
	}
}

func (b *backend) handleSettingAccount(c *gin.Context) {
	var req model.RequestSettingAccount
	if !requestBinding(c, &req) {
		return
	}

	if errDetail := b.Processor.SettingAccount(&req); errDetail != nil {
		errDetailLog(errDetail, b.SetLog, "Setting account failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.Status(http.StatusOK)
}
