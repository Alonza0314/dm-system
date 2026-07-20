package main_test

import (
	"backend/model"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/free-ran-ue/util"
)

var deviceRoutes = []route{
	newRoute("/device/cate", http.MethodGet),
	newRoute("/device/cate/dev", http.MethodGet),
	newRoute("/device", http.MethodPost),
	newRoute("/device/cate/dev", http.MethodDelete),
}

func TestApiDevice(t *testing.T) {
	testAuthRoutes(t, "Device", deviceRoutes)

	login(t)
	addCategory(t, category)

	t.Run("CreateDevice", testCreateDevice)
	t.Run("GetDevices", testGetDevices)
	t.Run("GetDevice", testGetDevice)
	t.Run("DeleteDevice", testDeleteDevice)
}

var (
	category = "cate"
	devices  = []string{
		"dev1",
		"dev2",
	}
)

func testCreateDevice(t *testing.T) {
	t.Run("Create 1 2", func(t *testing.T) {
		for _, dv := range devices {
			request := model.RequestCreateDevice{
				Category: category,
				Name:     dv,
			}

			requestByte, err := json.Marshal(request)
			if err != nil {
				handleJsonMarshalError(t, err)
			}

			response, err := util.SendHttpRequest(BASE_URL+"/device", http.MethodPost, header, requestByte)
			if err != nil {
				handleSendHttpError(t, err)
			}

			handleCheckStatusCode(t, http.StatusOK, response.StatusCode)
		}
	})

	t.Run("Duplicate create 1 2", func(t *testing.T) {
		for _, dv := range devices {
			request := model.RequestCreateDevice{
				Category: category,
				Name:     dv,
			}

			requestByte, err := json.Marshal(request)
			if err != nil {
				handleJsonMarshalError(t, err)
			}

			response, err := util.SendHttpRequest(BASE_URL+"/device", http.MethodPost, header, requestByte)
			if err != nil {
				handleSendHttpError(t, err)
			}

			handleCheckStatusCode(t, http.StatusConflict, response.StatusCode)
		}
	})

	t.Run("Non exist category", func(t *testing.T) {
		for _, dv := range devices {
			request := model.RequestCreateDevice{
				Category: "non-exist",
				Name:     dv,
			}

			requestByte, err := json.Marshal(request)
			if err != nil {
				handleJsonMarshalError(t, err)
			}

			response, err := util.SendHttpRequest(BASE_URL+"/device", http.MethodPost, header, requestByte)
			if err != nil {
				handleSendHttpError(t, err)
			}

			handleCheckStatusCode(t, http.StatusNotFound, response.StatusCode)
		}
	})
}

func testGetDevices(t *testing.T) {
	t.Run("Get all devices", func(t *testing.T) {
		response, err := util.SendHttpRequest(BASE_URL+"/device/"+category, http.MethodGet, header, nil)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

		var responseGetDevices model.ResponseGetDevices
		if err := json.Unmarshal(response.Body, &responseGetDevices); err != nil {
			handleJsonUnmarshalError(t, err)
		}

		if len(responseGetDevices.Devices) != len(devices) {
			t.Fatalf("failed to get devices with incorrect length, expected %d, got %d", len(responseGetDevices.Devices), len(categories))
		}

		for _, dv := range responseGetDevices.Devices {
			found := false

			for _, dvv := range devices {
				if dv.Name == dvv {
					found = true
					break
				}
			}

			if !found {
				t.Fatalf("could not find the device %s in testcase", dv.Name)
			}
		}

		for _, dvv := range devices {
			found := false

			for _, dv := range responseGetDevices.Devices {
				if dvv == dv.Name {
					found = true
					break
				}
			}

			if !found {
				t.Fatalf("could not find the device %s in response body", dvv)
			}
		}
	})

	t.Run("Get devices in non-exist category", func(t *testing.T) {
		response, err := util.SendHttpRequest(BASE_URL+"/device/non-exist", http.MethodGet, header, nil)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusNotFound, response.StatusCode)
	})
}

func testGetDevice(t *testing.T) {

}

func testDeleteDevice(t *testing.T) {

}
