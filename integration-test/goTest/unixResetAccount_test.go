package main_test

import (
	"backend/model"
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"testing"

	"github.com/free-ran-ue/util"
)

const (
	CONTAINER_NAME = "dm-test"
)

func TestUnixResetAccount(t *testing.T) {
	login(t)
	changeAccount(t)

	t.Run("ResetAccount", testResetAccount)
	t.Run("Relogin with default account", testReloginWithDefaultAccount)
}

func testResetAccount(t *testing.T) {
	cmd := exec.Command("docker", "exec", CONTAINER_NAME, "./tool", "-a", "resetAccount")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run reset account tool: %v\noutput: %s", err, output)
	}

	result := strings.ToLower(string(output))

	if strings.Contains(result, "fail") {
		t.Fatalf("reset account tool reported failure, output: %s", output)
	}

	if !strings.Contains(result, "success") {
		t.Fatalf("reset account tool output did not report success, output: %s", output)
	}
}

func testReloginWithDefaultAccount(t *testing.T) {
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
}
