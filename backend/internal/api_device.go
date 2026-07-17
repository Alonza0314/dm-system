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
			HandlerFunc: withLogging("GetDevices", b.DevLog, b.handleGetDevices),
		},
		{
			Name:        "GetDevice",
			Method:      http.MethodGet,
			Pattern:     "/device/:cate/:dev",
			HandlerFunc: withLogging("GetDevice", b.DevLog, b.handleGetDevice),
		},
		{
			Name:        "CreateDevice",
			Method:      http.MethodPost,
			Pattern:     "/device",
			HandlerFunc: withLogging("CreateDevice", b.DevLog, b.handleCreateDevice),
		},
		{
			Name:        "DeleteDevice",
			Method:      http.MethodDelete,
			Pattern:     "/device/:cate/:dev",
			HandlerFunc: withLogging("DeleteDevice", b.DevLog, b.handleDeleteDevice),
		},
	}
}

func (b *backend) handleGetDevices(c *gin.Context) {
	response, errDetail := b.Processor.GetDevices(c.Param("cate"))
	if errDetail != nil {
		errDetailLog(errDetail, b.DevLog, "Get devices failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (b *backend) handleGetDevice(c *gin.Context) {
	response, errDetail := b.Processor.GetDevice(c.Param("cate"), c.Param("dev"))
	if errDetail != nil {
		errDetailLog(errDetail, b.DevLog, "Get device %s failed for %s: %s", c.Param("dev"), c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (b *backend) handleCreateDevice(c *gin.Context) {
	var req model.RequestCreateDevice
	if !requestBinding(c, &req) {
		return
	}

	response, errDetail := b.Processor.CreateDevice(&req)
	if errDetail != nil {
		errDetailLog(errDetail, b.DevLog, "Create device failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (b *backend) handleDeleteDevice(c *gin.Context) {
	response, errDetail := b.Processor.DeleteDevice(c.Param("cate"), c.Param("dev"))
	if errDetail != nil {
		errDetailLog(errDetail, b.DevLog, "Delete device %s failed for %s: %s", c.Param("dev"), c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.JSON(http.StatusOK, response)
}
