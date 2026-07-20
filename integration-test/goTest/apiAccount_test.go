package main_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"backend/model"

	"github.com/free-ran-ue/util"
)

func TestApiAccount(t *testing.T) {
	t.Run("testLogin", testLogin)
	t.Run("testLogout", testLogout)
}

var testLoginCases = []struct {
	name     string
	username string
	password string

	expected
}{
	{
		name:     "login success",
		username: "admin",
		password: "0000",

		expected: expected{
			expectedStatusCode: http.StatusOK,
		},
	},
	{
		name:     "login unauthorized invalid username",
		username: "adminadmin",
		password: "0000",

		expected: expected{
			expectedStatusCode: http.StatusUnauthorized,
		},
	},
	{
		name:     "login unauthorized invalid password",
		username: "admin",
		password: "00000000",

		expected: expected{
			expectedStatusCode: http.StatusUnauthorized,
		},
	},
}

func testLogin(t *testing.T) {
	for _, tc := range testLoginCases {
		t.Run(tc.name, func(t *testing.T) {
			request := model.RequestLogin{
				Username: tc.username,
				Password: tc.password,
			}

			requestByte, err := json.Marshal(request)
			if err != nil {
				handleJsonMarshalError(t, err)
			}

			response, err := util.SendHttpRequest(BASE_URL+"/login", http.MethodPost, nil, requestByte)
			if err != nil {
				handleSendHttpError(t, err)
			}

			handleCheckdStatusCode(t, tc.expectedStatusCode, response.StatusCode)

			if response.StatusCode != http.StatusOK {
				return
			}

			var responseLogin model.ResponseLogin
			if err := json.Unmarshal(response.Body, &responseLogin); err != nil {
				handleJsonUnmarshalError(t, err)
			}

			if _, err := util.ValidateJWT(responseLogin.Token, JWT_SECRET); err != nil {
				t.Fatalf("invalid login jwt token: %v", err)
			}
		})
	}
}

func testLogout(t *testing.T) {
	response, err := util.SendHttpRequest(BASE_URL+"/logout", http.MethodPost, nil, nil)
	if err != nil {
		handleSendHttpError(t, err)
	}

	handleCheckdStatusCode(t, http.StatusNoContent, response.StatusCode)
}
