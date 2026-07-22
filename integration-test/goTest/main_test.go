package main_test

import (
	"backend/model"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/free-ran-ue/util"
)

const (
	BASE_URL   = "http://localhost:8888/api"
	JWT_SECRET = "dm-system"
)

var (
	header map[string]string

	category = "cate"
	device   = "dev"
	devices  = []string{
		"dev1",
		"dev2",
	}
	user = "tester"

	username1 = "username1"
	password1 = "password1"

	username2 = "username2"
	password2 = "password2"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()

	os.Exit(exitCode)
}

type expected struct {
	expectedStatusCode int
	expectedResponse   map[string]interface{}
}

func handleSendHttpError(t *testing.T, err error) {
	t.Fatalf("send http request failed: %v", err)
}

func handleJsonMarshalError(t *testing.T, err error) {
	t.Fatalf("json marshal error: %v", err)
}

func handleJsonUnmarshalError(t *testing.T, err error) {
	t.Fatalf("json unmarshal error: %v", err)
}

func handleCheckStatusCode(t *testing.T, expected, got int) {
	if expected == got {
		return
	}
	t.Fatalf("unexpected status code, expected: %d, got: %d", expected, got)
}

type route struct {
	route  string
	method string
}

func newRoute(rt, md string) route {
	return route{
		route:  rt,
		method: md,
	}
}

func testAuthRoutes(t *testing.T, name string, routes []route) {
	t.Run("Auth"+name, func(t *testing.T) {
		for _, rt := range routes {
			response, err := util.SendHttpRequest(BASE_URL+rt.route, rt.method, nil, nil)
			if err != nil {
				handleSendHttpError(t, err)
			}

			handleCheckStatusCode(t, http.StatusUnauthorized, response.StatusCode)
		}
	})
}

func login(t *testing.T) {
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

	handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("failed to login: %v", err)
	}

	var responseLogin model.ResponseLogin
	if err := json.Unmarshal(response.Body, &responseLogin); err != nil {
		handleJsonUnmarshalError(t, err)
	}

	header = make(map[string]string, 0)
	header["Authorization"] = "Bearer " + responseLogin.Token
}

func addCategory(t *testing.T, cate string) {
	request := model.RequestCreateCategory{
		Name: cate,
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

func addDevice(t *testing.T, dev string) {
	request := model.RequestCreateDevice{
		Category: category,
		Name:     dev,
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
