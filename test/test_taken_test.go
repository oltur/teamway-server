package test

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/oltur/teamway-server/controller"
	"github.com/oltur/teamway-server/model"
	"net/http"
	"net/http/httptest"
	neturl "net/url"
	"testing"
)

func TestTakeTestOk(t *testing.T) {
	model.InitModel()
	test, err := getFirstTest()
	if err != nil {
		t.Fatal(err)
	}
	router, c := controller.SetupRouter()
	gwtToken, _, _, err := c.DoLogin("User1", "1")
	if err != nil {
		t.Fatal(err)
	}

	result, err := doTakeTestOk(router, test, gwtToken)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, result, test.NegativeResult)
}

func TestTakeTestFailWrongAnswer(t *testing.T) {
	model.InitModel()
	test, err := getFirstTest()
	if err != nil {
		t.Fatal(err)
	}
	router, c := controller.SetupRouter()
	gwtToken, _, _, err := c.DoLogin("User1", "1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = doTakeTestFailWrongAnswer(router, test, gwtToken)
	if err == nil {
		err = fmt.Errorf("unexpected success")
		t.Fatal(err)
	}
	if err.Error() != "unexpected http code" {
		err = fmt.Errorf("unexpected error message")
		t.Fatal(err)
	}
}

// ------- implementation details ---------------

func doTakeTestOk(router *gin.Engine, test *model.Test, gwtToken string) (result string, err error) {
	// for all questions
	for {
		// get next question
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/api/v1/test-taken/next?test-id=%s", string(test.ID))
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+gwtToken)
		router.ServeHTTP(w, req)
		if w.Code == 204 {
			break
		}
		if w.Code != 200 {
			err = fmt.Errorf("unexpected http code")
			return
		}
		body := w.Body.String()
		var nextQuestion string
		err = json.Unmarshal([]byte(body), &nextQuestion)
		if err != nil {
			return
		}

		if nextQuestion == "" {
			break
		}

		firstAnswer := findQuestion(test.Questions, nextQuestion).Answers[0]

		// check interim result (202)
		w = httptest.NewRecorder()
		url = fmt.Sprintf("/api/v1/test-taken?test-id=%s", string(test.ID))
		req, _ = http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+gwtToken)
		router.ServeHTTP(w, req)
		if w.Code != 202 {
			err = fmt.Errorf("unexpected http code")
			return
		}

		// answer
		w = httptest.NewRecorder()
		url = fmt.Sprintf("/api/v1/test-taken?test-id=%s&question-title=%s&answer-title=%s",
			string(test.ID), neturl.QueryEscape(nextQuestion), neturl.QueryEscape(firstAnswer.Title))
		req, _ = http.NewRequest("POST", url, nil)
		req.Header.Add("Authorization", "Bearer "+gwtToken)
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			err = fmt.Errorf("unexpected http code")
			return
		}
	}
	// check result
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/v1/test-taken?test-id=%s", string(test.ID))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		err = fmt.Errorf("unexpected http code")
		return
	}
	body := w.Body.String()
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return
	}

	if result != test.NegativeResult {
		err = fmt.Errorf("unexpected test result")
		return
	}

	return
}

func doTakeTestFailWrongAnswer(router *gin.Engine, test *model.Test, gwtToken string) (result string, err error) {
	// for all questions
	for {
		// get next question
		w := httptest.NewRecorder()
		url := fmt.Sprintf("/api/v1/test-taken/next?test-id=%s", string(test.ID))
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+gwtToken)
		router.ServeHTTP(w, req)
		if w.Code == 204 {
			break
		}
		if w.Code != 200 {
			err = fmt.Errorf("unexpected http code")
			return
		}
		body := w.Body.String()
		var nextQuestion string
		err = json.Unmarshal([]byte(body), &nextQuestion)
		if err != nil {
			return
		}

		if nextQuestion == "" {
			break
		}

		firstAnswer := findQuestion(test.Questions, nextQuestion).Answers[0]

		// check interim result (202)
		w = httptest.NewRecorder()
		url = fmt.Sprintf("/api/v1/test-taken?test-id=%s", string(test.ID))
		req, _ = http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+gwtToken)
		router.ServeHTTP(w, req)
		if w.Code != 202 {
			err = fmt.Errorf("unexpected http code")
			return
		}

		// answer
		w = httptest.NewRecorder()
		url = fmt.Sprintf("/api/v1/test-taken?test-id=%s&question-title=%s&answer-title=%s",
			string(test.ID), neturl.QueryEscape(nextQuestion), neturl.QueryEscape(firstAnswer.Title+" !!!!! WRONG !!!!!"))
		req, _ = http.NewRequest("POST", url, nil)
		req.Header.Add("Authorization", "Bearer "+gwtToken)
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			err = fmt.Errorf("unexpected http code")
			return
		}
	}
	// check result
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/v1/test-taken?test-id=%s", string(test.ID))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+gwtToken)
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		err = fmt.Errorf("unexpected http code")
		return
	}
	body := w.Body.String()
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return
	}

	if result != test.NegativeResult {
		err = fmt.Errorf("unexpected test result")
		return
	}

	return
}
