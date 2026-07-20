package main_test

import (
	"os"
	"testing"
)

const (
	BASE_URL   = "http://localhost:8888/api"
	JWT_SECRET = "dm-system"
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

func handleCheckdStatusCode(t *testing.T, expected, got int) {
	if expected == got {
		return
	}
	t.Fatalf("unexpected status code, expected: %d, got: %d", expected, got)
}
