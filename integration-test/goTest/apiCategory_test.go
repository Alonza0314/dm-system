package main_test

import (
	"backend/model"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/free-ran-ue/util"
)

var categoryRoutes = []route{
	newRoute("/category", http.MethodGet),
	newRoute("/category/cate", http.MethodGet),
	newRoute("/category", http.MethodPost),
	newRoute("/category/cate", http.MethodDelete),
}

func TestApiCategory(t *testing.T) {
	testAuthRoutes(t, "Category", categoryRoutes)

	login(t)

	t.Run("CreateCategory", testCreateCategory)
	t.Run("GetCategories", testGetCategories)
	t.Run("GetCategorey", testGetCategory)
	t.Run("DeleteCategory", testDeleteCategory)
}

var categories = []string{
	"cate1",
	"cate2",
}

func testCreateCategory(t *testing.T) {
	t.Run("Create 1 2", func(t *testing.T) {
		for _, ct := range categories {
			request := model.RequestCreateCategory{
				Name: ct,
			}
			requestByte, err := json.Marshal(request)
			if err != nil {
				handleJsonMarshalError(t, err)
			}

			response, err := util.SendHttpRequest(BASE_URL+"/category", http.MethodPost, header, requestByte)
			if err != nil {
				handleSendHttpError(t, err)
			}

			handleCheckStatusCode(t, http.StatusCreated, response.StatusCode)
		}
	})

	t.Run("Duplicate create 1 2", func(t *testing.T) {
		for _, ct := range categories {
			request := model.RequestCreateCategory{
				Name: ct,
			}
			requestByte, err := json.Marshal(request)
			if err != nil {
				handleJsonMarshalError(t, err)
			}

			response, err := util.SendHttpRequest(BASE_URL+"/category", http.MethodPost, header, requestByte)
			if err != nil {
				handleSendHttpError(t, err)
			}

			handleCheckStatusCode(t, http.StatusConflict, response.StatusCode)
		}
	})
}

func testGetCategories(t *testing.T) {
	response, err := util.SendHttpRequest(BASE_URL+"/category", http.MethodGet, header, nil)
	if err != nil {
		handleSendHttpError(t, err)
	}

	handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

	var responseGetCategories model.ResponseGetCategories
	if err := json.Unmarshal(response.Body, &responseGetCategories); err != nil {
		handleJsonUnmarshalError(t, err)
	}

	if len(responseGetCategories.Categories) != len(categories) {
		t.Fatalf("failed to get categories with incorrect length, expected %d, got %d", len(responseGetCategories.Categories), len(categories))
	}

	for _, ct := range responseGetCategories.Categories {
		found := false

		for _, cct := range categories {
			if ct.Name == cct {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("could not find the category %s in testcase", ct.Name)
		}
	}

	for _, cct := range categories {
		found := false

		for _, ct := range responseGetCategories.Categories {
			if cct == ct.Name {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("could not find the category %s in response body", cct)
		}
	}
}

func testGetCategory(t *testing.T) {
	for _, ct := range categories {
		t.Run("Get category "+ct, func(t *testing.T) {
			response, err := util.SendHttpRequest(BASE_URL+"/category/"+ct, http.MethodGet, header, nil)
			if err != nil {
				handleSendHttpError(t, err)
			}

			handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

			var responseGetCategory model.ResponseGetCategory
			if err := json.Unmarshal(response.Body, &responseGetCategory); err != nil {
				handleJsonUnmarshalError(t, err)
			}

			if responseGetCategory.Name != ct {
				t.Fatalf("incorrect response category name, expected %s, got %s", ct, responseGetCategory.Name)
			}

			if responseGetCategory.IdleDevice != 0 {
				t.Fatalf("init idle device is not 0 in category %s", ct)
			}

			if responseGetCategory.UsingDevice != 0 {
				t.Fatalf("init using device is not 0 in category %s", ct)
			}
		})
	}
}

func testDeleteCategory(t *testing.T) {
	t.Run("Delete 1 2", func(t *testing.T) {
		for _, ct := range categories {
			response, err := util.SendHttpRequest(BASE_URL+"/category/"+ct, http.MethodDelete, header, nil)
			if err != nil {
				handleSendHttpError(t, err)
			}

			handleCheckStatusCode(t, http.StatusOK, response.StatusCode)
		}
	})

	t.Run("Duplicate delete 1 2", func(t *testing.T) {
		for _, ct := range categories {
			response, err := util.SendHttpRequest(BASE_URL+"/category/"+ct, http.MethodDelete, header, nil)
			if err != nil {
				handleSendHttpError(t, err)
			}

			handleCheckStatusCode(t, http.StatusNotFound, response.StatusCode)
		}
	})

	t.Run("Check category is empty", func(t *testing.T) {
		response, err := util.SendHttpRequest(BASE_URL+"/category", http.MethodGet, header, nil)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

		var responseGetCategories model.ResponseGetCategories
		if err := json.Unmarshal(response.Body, &responseGetCategories); err != nil {
			handleJsonUnmarshalError(t, err)
		}

		if len(responseGetCategories.Categories) != 0 {
			t.Fatalf("failed to get categories with incorrect length, expected %d, got %d", len(responseGetCategories.Categories), 0)
		}
	})
}
