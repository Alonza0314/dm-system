package main_test

import (
	"backend/model"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/free-ran-ue/util"
)

var historyRoutes = []route{
	newRoute("/history/cate/dev", http.MethodGet),
}

func TestApiHistory(t *testing.T) {
	testAuthRoutes(t, "History", historyRoutes)

	login(t)
	addCategory(t, category)
	addDevice(t, device)

	testGetHistory(t)
	testMaxHistory(t)
}

func testGetHistory(t *testing.T) {
	t.Run("No record", func(t *testing.T) {
		response, err := util.SendHttpRequest(BASE_URL+"/history/"+category+"/"+device, http.MethodGet, header, nil)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

		var history model.ResponseGetHistory
		if err := json.Unmarshal(response.Body, &history); err != nil {
			handleJsonUnmarshalError(t, err)
		}

		if len(history.Histories) != 0 {
			t.Fatalf("failed to get history which should be 0 record, but got %d records", len(history.Histories))
		}
	})

	borrow(t, category, device, user)
	returnn(t, category, device, user)

	t.Run("One record", func(t *testing.T) {
		response, err := util.SendHttpRequest(BASE_URL+"/history/"+category+"/"+device, http.MethodGet, header, nil)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

		var history model.ResponseGetHistory
		if err := json.Unmarshal(response.Body, &history); err != nil {
			handleJsonUnmarshalError(t, err)
		}

		if len(history.Histories) != 1 {
			t.Fatalf("failed to get history which should be 1 record, but got %d records", len(history.Histories))
		}

		if history.Histories[0].User != user {
			t.Fatalf("failed to get history with incorrect user, expected %s, but got %s", user, history.Histories[0].User)
		}
	})
}

func testMaxHistory(t *testing.T) {
	for i := 1; i <= 20; i += 1 {
		borrow(t, category, device, fmt.Sprintf("%s-%d", user, i))
		returnn(t, category, device, fmt.Sprintf("%s-%d", user, i))
	}

	t.Run("Max records", func(t *testing.T) {
		response, err := util.SendHttpRequest(BASE_URL+"/history/"+category+"/"+device, http.MethodGet, header, nil)
		if err != nil {
			handleSendHttpError(t, err)
		}

		handleCheckStatusCode(t, http.StatusOK, response.StatusCode)

		var history model.ResponseGetHistory
		if err := json.Unmarshal(response.Body, &history); err != nil {
			handleJsonUnmarshalError(t, err)
		}

		if len(history.Histories) != 10 {
			t.Fatalf("failed to get history which should be 10 record, but got %d records", len(history.Histories))
		}

		for i, j := 20, 0; i > 10; i, j = i-1, j+1 {
			if history.Histories[j].User != fmt.Sprintf("%s-%d", user, i) {
				t.Fatalf("failed to verify the history with username, expected %s, but got %s", fmt.Sprintf("%s-%d", user, i), history.Histories[j].User)
			}
		}
	})
}
