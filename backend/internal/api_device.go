package internal

import (
	"backend/model"
	"net/http"

	"github.com/free-ran-ue/util"
	"github.com/gin-gonic/gin"
)

func (b *backend) getDeviceRoutes() util.Routes {
	return util.Routes{
		{
			Name:        "GetDevices",
			Method:      http.MethodGet,
			Pattern:     "/device/:cate",
			HandlerFunc: b.handleGetDevices,
		},
		{
			Name:        "GetDevice",
			Method:      http.MethodGet,
			Pattern:     "/device/:cate/:dev",
			HandlerFunc: b.handleGetDevice,
		},
		{
			Name:        "CreateDevice",
			Method:      http.MethodPost,
			Pattern:     "/device",
			HandlerFunc: b.handleCreateDevice,
		},
		{
			Name:        "DeleteDevice",
			Method:      http.MethodDelete,
			Pattern:     "/device/:cate/:dev",
			HandlerFunc: b.handleDeleteDevice,
		},
	}
}

func (b *backend) handleGetDevices(c *gin.Context) {
	b.DevLog.Infof("Get devices request from %s", c.ClientIP())

	response, errDetail := b.Processor.GetDevices(c.Param("cate"))
	if errDetail != nil {
		errDetailLog(errDetail, b.DevLog, "Get devices failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	b.DevLog.Infof("Get devices successful from %s", c.ClientIP())
	c.JSON(http.StatusOK, response)
}

func (b *backend) handleGetDevice(c *gin.Context) {
	b.DevLog.Infof("Get device request from %s", c.ClientIP())

	response, errDetail := b.Processor.GetDevice(c.Param("cate"), c.Param("dev"))
	if errDetail != nil {
		errDetailLog(errDetail, b.CatLog, "Get device %s failed for %s: %s", c.Param("dev"), c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	b.DevLog.Infof("Get device successful from %s", c.ClientIP())
	c.JSON(http.StatusOK, response)
}

func (b *backend) handleCreateDevice(c *gin.Context) {
	b.DevLog.Infof("Create device request from %s", c.ClientIP())

	var req model.RequestCreateDevice
	if !requestBinding(c, &req) {
		b.DevLog.Warnf("Invalid create device request from %s", c.ClientIP())
		return
	}

	response, errDetail := b.Processor.CreateDevice(&req)
	if errDetail != nil {
		errDetailLog(errDetail, b.CatLog, "Create device failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	b.DevLog.Infof("Create device successful from %s", c.ClientIP())
	c.JSON(http.StatusOK, response)
}

func (b *backend) handleDeleteDevice(c *gin.Context) {
	b.DevLog.Infof("Delete device request from %s", c.ClientIP())

	response, errDetail := b.Processor.DeleteDevice(c.Param("cate"), c.Param("dev"))
	if errDetail != nil {
		errDetailLog(errDetail, b.CatLog, "Delete Device %s failed for %s: %s", c.Param("dev"), c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	b.DevLog.Infof("Delete device successful from %s", c.ClientIP())
	c.JSON(http.StatusOK, response)
}
