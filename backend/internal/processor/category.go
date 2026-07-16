package processor

import (
	"backend/constant"
	"backend/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func (p *Processor) GetCategories() (*model.ResponseGetCategories, *model.ErrorDetail) {
	categories, err := p.DmContext.Db().LoadAll(constant.COLL_CATEGORY)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get categories: %v", err),
		}
	}

	response := &model.ResponseGetCategories{
		Categories: make([]model.Category, 0, len(categories)),
	}
	for _, v := range categories {
		var category model.Category
		if err := json.Unmarshal([]byte(v), &category); err != nil {
			return nil, &model.ErrorDetail{
				HttpStatus: http.StatusInternalServerError,
				Detail:     fmt.Sprintf("failed to unmarshal categories: %v", err),
			}
		}

		response.Categories = append(response.Categories, category)
	}

	return response, nil
}

func (p *Processor) GetCategory(cate string) (*model.ResponseGetCategory, *model.ErrorDetail) {
	exist, err := p.DmContext.Exist(constant.COLL_CATEGORY, cate)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get category: %v", err),
		}
	}
	if !exist {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusNotFound,
			Detail:     fmt.Sprintf("category %s not found", cate),
		}
	}

	category, err := p.DmContext.Db().Load(constant.COLL_CATEGORY, cate)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get category: %v", err),
		}
	}

	var categoryUnmarshal model.ResponseGetCategory
	if err := json.Unmarshal([]byte(category), &categoryUnmarshal); err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to unmarshal category: %v", err),
		}
	}

	return &categoryUnmarshal, nil
}

func (p *Processor) CreateCategory(req *model.RequestCreateCategory) (*model.ResponseCreateCategory, *model.ErrorDetail) {
	exist, err := p.DmContext.Db().Exist(constant.COLL_CATEGORY, req.Name)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to check if category exists: %v", err),
		}
	}
	if exist {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusConflict,
			Detail:     "category already exists",
		}
	}

	catId, err := p.DmContext.RequestId(constant.ID_KEY_CATEGORY)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to request category id: %v", err),
		}
	}
	req.Id = catId

	category, err := json.Marshal(req)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprint("failed to marshal json: %v", err),
		}
	}

	if err := p.DmContext.Db().Save(constant.COLL_CATEGORY, req.Name, string(category)); err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to save category: %v", err),
		}
	}

	return &model.ResponseCreateCategory{
		Message: "Create successful",
	}, nil
}

func (p *Processor) DeleteCategory(cate string) (*model.ResponseDeleteCategory, *model.ErrorDetail) {
	exist, err := p.DmContext.Exist(constant.COLL_CATEGORY, cate)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get category: %v", err),
		}
	}
	if !exist {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusNotFound,
			Detail:     fmt.Sprintf("category %s not found", cate),
		}
	}

	devices, err := p.DmContext.LoadAll(constant.COLL_CATEGORY_TAG + cate)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get category devices: %v", err),
		}
	}

	if len(devices) != 0 {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusConflict,
			Detail:     fmt.Sprintf("still exist devices under category %s", cate),
		}
	}

	if err := p.DmContext.Remove(constant.COLL_CATEGORY, cate); err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to delete category: %v", err),
		}
	}

	return &model.ResponseDeleteCategory{
		Message: "Delete successful",
	}, nil
}
