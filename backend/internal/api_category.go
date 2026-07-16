package internal

import (
	"backend/model"
	"net/http"

	"github.com/free-ran-ue/util"
	"github.com/gin-gonic/gin"
)

func (b *backend) getCategoryRoutes() util.Routes {
	return util.Routes{
		{
			Name:        "GetCategories",
			Method:      http.MethodGet,
			Pattern:     "/category",
			HandlerFunc: b.handleGetCategories,
		},
		{
			Name:        "GetCategory",
			Method:      http.MethodGet,
			Pattern:     "/category/:cate",
			HandlerFunc: b.handleGetCategory,
		},
		{
			Name:        "CreateCategory",
			Method:      http.MethodPost,
			Pattern:     "/category",
			HandlerFunc: b.handleCreateCategory,
		},
		{
			Name:        "DeleteCategory",
			Method:      http.MethodDelete,
			Pattern:     "/category/:cate",
			HandlerFunc: b.handleDeleteCategory,
		},
	}
}

func (b *backend) handleGetCategories(c *gin.Context) {
	b.CatLog.Infof("Get categories request from %s", c.ClientIP())

	response, errDetail := b.Processor.GetCategories()
	if errDetail != nil {
		errDetailLog(errDetail, b.CatLog, "Get categories failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	b.CatLog.Infof("Get categories successful from %s", c.ClientIP())
	c.JSON(http.StatusOK, response)
}

func (b *backend) handleGetCategory(c *gin.Context) {
	b.CatLog.Infof("Get category request from %s", c.ClientIP())

	response, errDetail := b.Processor.GetCategory(c.Param("cate"))
	if errDetail != nil {
		errDetailLog(errDetail, b.CatLog, "Get category %s failed for %s: %s", c.Param("cat"), c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	b.CatLog.Infof("Get category successful from %s", c.ClientIP())
	c.JSON(http.StatusOK, response)
}

func (b *backend) handleCreateCategory(c *gin.Context) {
	b.CatLog.Infof("Create category request from %s", c.ClientIP())

	var req model.RequestCreateCategory
	if !requestBinding(c, &req) {
		b.CatLog.Warnf("Invalid create category request from %s", c.ClientIP())
		return
	}

	response, errDetail := b.Processor.CreateCategory(&req)
	if errDetail != nil {
		errDetailLog(errDetail, b.CatLog, "Create category failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	b.CatLog.Infof("Create category successful from %s", c.ClientIP())
	c.JSON(http.StatusCreated, response)
}

func (b *backend) handleDeleteCategory(c *gin.Context) {
	b.CatLog.Infof("Delete category request from %s", c.ClientIP())

	response, errDetail := b.Processor.DeleteCategory(c.Param("cate"))
	if errDetail != nil {
		errDetailLog(errDetail, b.CatLog, "Delete category %s failed for %s: %s", c.Param("cat"), c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	b.CatLog.Infof("Delete category successful from %s", c.ClientIP())
	c.JSON(http.StatusOK, response)
}
