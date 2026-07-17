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
			HandlerFunc: withLogging("GetCategories", b.CatLog, b.handleGetCategories),
		},
		{
			Name:        "GetCategory",
			Method:      http.MethodGet,
			Pattern:     "/category/:cate",
			HandlerFunc: withLogging("GetCategory", b.CatLog, b.handleGetCategory),
		},
		{
			Name:        "CreateCategory",
			Method:      http.MethodPost,
			Pattern:     "/category",
			HandlerFunc: withLogging("CreateCategory", b.CatLog, b.handleCreateCategory),
		},
		{
			Name:        "DeleteCategory",
			Method:      http.MethodDelete,
			Pattern:     "/category/:cate",
			HandlerFunc: withLogging("DeleteCategory", b.CatLog, b.handleDeleteCategory),
		},
	}
}

func (b *backend) handleGetCategories(c *gin.Context) {
	response, errDetail := b.Processor.GetCategories()
	if errDetail != nil {
		errDetailLog(errDetail, b.CatLog, "Get categories failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (b *backend) handleGetCategory(c *gin.Context) {
	response, errDetail := b.Processor.GetCategory(c.Param("cate"))
	if errDetail != nil {
		errDetailLog(errDetail, b.CatLog, "Get category %s failed for %s: %s", c.Param("cate"), c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (b *backend) handleCreateCategory(c *gin.Context) {
	var req model.RequestCreateCategory
	if !requestBinding(c, &req) {
		return
	}

	response, errDetail := b.Processor.CreateCategory(&req)
	if errDetail != nil {
		errDetailLog(errDetail, b.CatLog, "Create category failed for %s: %s", c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (b *backend) handleDeleteCategory(c *gin.Context) {
	response, errDetail := b.Processor.DeleteCategory(c.Param("cate"))
	if errDetail != nil {
		errDetailLog(errDetail, b.CatLog, "Delete category %s failed for %s: %s", c.Param("cate"), c.ClientIP(), errDetail.Detail)
		c.JSON(errDetail.HttpStatus, errDetail)
		return
	}

	c.JSON(http.StatusOK, response)
}
