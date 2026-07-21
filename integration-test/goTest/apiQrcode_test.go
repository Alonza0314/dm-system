package main_test

import (
	"backend/constant"
	"backend/model"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/free-ran-ue/util"
)

func TestApiQrcode(t *testing.T) {
	login(t)
	addCategory(t, category)
	addDevice(t, device)

	t.Run("Borrow", testBorrow)
	t.Run("Return", testReturn)
}

func testBorrow(t *testing.T) {
	t.Run("Non exist category", func(t *testing.T) {
		request := &model.RequestQrcodeBorrow{
			User: user,
		}
		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}

		response, err := util.SendHttpRequest(BASE_URL+"/qrcode/non-exist/non-exist", http.MethodPost, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusNotFound, response.StatusCode)
	})

	t.Run("Non exist device", func(t *testing.T) {
		request := &model.RequestQrcodeBorrow{
			User: user,
		}
		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}

		response, err := util.SendHttpRequest(BASE_URL+"/qrcode/"+category+"/non-exist", http.MethodPost, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusNotFound, response.StatusCode)
	})

	t.Run("Borrow", func(t *testing.T) {
		request := &model.RequestQrcodeBorrow{
			User: user,
		}
		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}

		response, err := util.SendHttpRequest(BASE_URL+"/qrcode/"+category+"/"+device, http.MethodPost, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)
	})

	t.Run("Duplicate borrow", func(t *testing.T) {
		request := &model.RequestQrcodeBorrow{
			User: user,
		}
		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}

		response, err := util.SendHttpRequest(BASE_URL+"/qrcode/"+category+"/"+device, http.MethodPost, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusConflict, response.StatusCode)
	})

	t.Run("Check category using device", func(t *testing.T) {
		response, err := util.SendHttpRequest(BASE_URL+"/category/"+category, http.MethodGet, header, nil)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

		var responseGetCategory model.ResponseGetCategory
		if err := json.Unmarshal(response.Body, &responseGetCategory); err != nil {
			handleJsonUnmarshalError(t, err)
		}

		if responseGetCategory.IdleDevice != 0 {
			t.Fatalf("idle device is not 0 in category %s", category)
		}

		if responseGetCategory.UsingDevice == 0 {
			t.Fatalf("using device is 0 in category %s", category)
		}
	})

	t.Run("Check device status", func(t *testing.T) {
		response, err := util.SendHttpRequest(BASE_URL+"/device/"+category+"/"+device, http.MethodGet, header, nil)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

		var responseGetDevice model.ResponseGetDevice
		if err := json.Unmarshal(response.Body, &responseGetDevice); err != nil {
			handleJsonUnmarshalError(t, err)
		}

		if responseGetDevice.Status == constant.STATUS_IDLE {
			t.Fatalf("incorrect device %s status, expected %s, got %s", device, constant.STATUS_USING, responseGetDevice.Status)
		}

		if responseGetDevice.User != user {
			t.Fatalf("incorrect device %s user, expected %s, got %s", device, user, responseGetDevice.User)
		}
	})
}

func testReturn(t *testing.T) {
	t.Run("Non exist category", func(t *testing.T) {
		request := &model.RequestQrcodeBorrow{
			User: user,
		}
		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}

		response, err := util.SendHttpRequest(BASE_URL+"/qrcode/non-exist/non-exist", http.MethodDelete, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusNotFound, response.StatusCode)
	})

	t.Run("Non exist device", func(t *testing.T) {
		request := &model.RequestQrcodeBorrow{
			User: user,
		}
		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}

		response, err := util.SendHttpRequest(BASE_URL+"/qrcode/"+category+"/non-exist", http.MethodDelete, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusNotFound, response.StatusCode)
	})

	t.Run("Return", func(t *testing.T) {
		request := &model.RequestQrcodeBorrow{
			User: user,
		}
		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}

		response, err := util.SendHttpRequest(BASE_URL+"/qrcode/"+category+"/"+device, http.MethodDelete, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)
	})

	t.Run("Duplicate return", func(t *testing.T) {
		request := &model.RequestQrcodeBorrow{
			User: user,
		}
		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}

		response, err := util.SendHttpRequest(BASE_URL+"/qrcode/"+category+"/"+device, http.MethodDelete, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusConflict, response.StatusCode)
	})

	t.Run("Check category idle device", func(t *testing.T) {
		response, err := util.SendHttpRequest(BASE_URL+"/category/"+category, http.MethodGet, header, nil)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

		var responseGetCategory model.ResponseGetCategory
		if err := json.Unmarshal(response.Body, &responseGetCategory); err != nil {
			handleJsonUnmarshalError(t, err)
		}

		if responseGetCategory.IdleDevice == 0 {
			t.Fatalf("idle device is 0 in category %s", category)
		}

		if responseGetCategory.UsingDevice != 0 {
			t.Fatalf("using device is not 0 in category %s", category)
		}
	})

	t.Run("Check device status", func(t *testing.T) {
		response, err := util.SendHttpRequest(BASE_URL+"/device/"+category+"/"+device, http.MethodGet, header, nil)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

		var responseGetDevice model.ResponseGetDevice
		if err := json.Unmarshal(response.Body, &responseGetDevice); err != nil {
			handleJsonUnmarshalError(t, err)
		}

		if responseGetDevice.Status == constant.STATUS_USING {
			t.Fatalf("incorrect device %s status, expected %s, got %s", device, constant.STATUS_IDLE, responseGetDevice.Status)
		}

		if responseGetDevice.User != "" {
			t.Fatalf("incorrect device %s user, expected nil, got %s", device, responseGetDevice.User)
		}
	})
}
