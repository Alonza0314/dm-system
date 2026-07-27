package main_test

import (
	"backend/model"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/free-ran-ue/util"
)

var settingRoutes = []route{
	newRoute("/setting/account", http.MethodPost),
}

func TestApiSetting(t *testing.T) {
	testAuthRoutes(t, "Setting", settingRoutes)

	login(t)

	t.Run("SettingAccount", testSettingAccount)
}

func testSettingAccount(t *testing.T) {
	t.Run("Update account with username 1", func(t *testing.T) {
		request := model.RequestSettingAccount{
			Username: username1,
			Password: password1,
		}

		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}

		response, err := util.SendHttpRequest(BASE_URL+"/setting/account", http.MethodPost, header, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)
	})

	t.Run("Login with old account should fail", func(t *testing.T) {
		request := model.RequestLogin{
			Username: "admin",
			Password: "0000",
		}

		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}
		response, err := util.SendHttpRequest(BASE_URL+"/login", http.MethodPost, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("Login with username 1 should success", func(t *testing.T) {
		request := model.RequestLogin{
			Username: username1,
			Password: password1,
		}

		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}
		response, err := util.SendHttpRequest(BASE_URL+"/login", http.MethodPost, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

		var responseLogin model.ResponseLogin
		if err := json.Unmarshal(response.Body, &responseLogin); err != nil {
			handleJsonUnmarshalError(t, err)
		}

		header = make(map[string]string, 0)
		header["Authorization"] = "Bearer " + responseLogin.Token
	})

	t.Run("Update account with username 2", func(t *testing.T) {
		request := model.RequestSettingAccount{
			Username: username2,
			Password: password2,
		}

		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}

		response, err := util.SendHttpRequest(BASE_URL+"/setting/account", http.MethodPost, header, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)
	})

	t.Run("Login with old account should fail", func(t *testing.T) {
		request := model.RequestLogin{
			Username: "admin",
			Password: "0000",
		}

		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}
		response, err := util.SendHttpRequest(BASE_URL+"/login", http.MethodPost, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("Login with username 1 should fail", func(t *testing.T) {
		request := model.RequestLogin{
			Username: username1,
			Password: password1,
		}

		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}
		response, err := util.SendHttpRequest(BASE_URL+"/login", http.MethodPost, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("Login with username 2 should success", func(t *testing.T) {
		request := model.RequestLogin{
			Username: username2,
			Password: password2,
		}

		requestByte, err := json.Marshal(request)
		if err != nil {
			handleJsonMarshalError(t, err)
		}
		response, err := util.SendHttpRequest(BASE_URL+"/login", http.MethodPost, nil, requestByte)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)
	})
}
