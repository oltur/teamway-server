package test

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/oltur/teamway-server/controller"
	"github.com/oltur/teamway-server/httputil"
	"github.com/oltur/teamway-server/model"
)

func TestGetAllTestsOk(t *testing.T) {
	model.InitModel()
	router, _ := controller.SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	var data []*model.Test
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 2 {
		t.Fatal("not enough tests")
	}
}

func TestGetTestsOk(t *testing.T) {
	model.InitModel()
	router, _ := controller.SetupRouter()

	test0, err := getFirstTest()
	if err != nil {
		t.Fatal(err)
	}
	testId := test0.ID

	// get test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/test/%s", testId), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	body := w.Body.String()
	var test *model.Test
	err = json.Unmarshal([]byte(body), &test)
	if err != nil {
		t.Fatal(err)
	}
	if test == nil {
		t.Fatal("no test")
	}
	if testId != test.ID {
		t.Fatal("not matching id")
	}
	if len(test.Questions) == 0 {
		t.Fatal("not enough questions")
	}
	question := test.Questions[0]
	if question.Title != "Question #1" {
		t.Fatal("wrong question")
	}
	if len(question.Answers) == 0 {
		t.Fatal("not enough answers")
	}
	answer := question.Answers[0]
	if answer.Title != "1" {
		t.Fatal("wrong answers")
	}
	if answer.Score != 1 {
		t.Fatal("wrong score")
	}
}

func TestGetTestFailedDoesNotExist(t *testing.T) {
	model.InitModel()
	router, _ := controller.SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/test/%s", xid.New().String()), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	body := w.Body.String()
	var data httputil.HTTPError
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	if data.Message != model.ErrNotFound.Error() {
		t.Fatal("not the right message")
	}
}
