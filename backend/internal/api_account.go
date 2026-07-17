package internal

import (
	"backend/model"
	"net/http"

	"github.com/free-ran-ue/util"
	"github.com/gin-gonic/gin"
)

func (b *backend) getAccountRoutes() util.Routes {
	return util.Routes{
		{
			Name:        "Login",
			Method:      http.MethodPost,
			Pattern:     "/login",
			HandlerFunc: withLogging("Login", b.AccLog, b.handleLogin),
		},
		{
			Name:        "Logout",
			Method:      http.MethodPost,
			Pattern:     "/logout",
			HandlerFunc: withLogging("Logout", b.AccLog, b.handleLogout),
		},
	}
}

func (b *backend) handleLogin(c *gin.Context) {
	var req model.RequestLogin
	if !requestBinding(c, &req) {
		return
	}

	response, errDetail := b.Processor.Login(&req)
	if errDetail != nil {
		errDetailLog(errDetail, b.AccLog, "Login failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (b *backend) handleLogout(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
