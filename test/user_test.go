package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/oltur/teamway-server/controller"
	"github.com/oltur/teamway-server/httputil"
	"github.com/oltur/teamway-server/model"
)

func TestGetUserOk(t *testing.T) {
	model.InitModel()
	router, c := controller.SetupRouter()
	w := httptest.NewRecorder()

	gwtToken, _, userId, err := c.DoLogin("User1", "1")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/user/%s", userId), nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	var data model.User
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	if data.UserName != "User1" {
		t.Fatal("not the right user")
	}
}

func TestGetUserFailedDoesNotExist(t *testing.T) {
	model.InitModel()
	router, c := controller.SetupRouter()
	w := httptest.NewRecorder()

	gwtToken, _, _, err := c.DoLogin("User1", "1")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "/api/v1/user/999", nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)

	body := w.Body.String()
	var data httputil.HTTPError
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	if data.Message != model.ErrAccessDenied.Error() {
		t.Fatal("not the right message")
	}
}

func TestGetUserFailedAccessDeniedNotAdmin(t *testing.T) {
	model.InitModel()
	router, c := controller.SetupRouter()
	w := httptest.NewRecorder()

	gwtToken, _, _, err := c.DoLogin("User1", "1")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", "/api/v1/user/2", nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)

	body := w.Body.String()
	var data httputil.HTTPError
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	if data.Message != model.ErrAccessDenied.Error() {
		t.Fatal("not the right message")
	}
}
