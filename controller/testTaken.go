package controller

import (
	"github.com/oltur/teamway-server/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oltur/teamway-server/httputil"
	"github.com/oltur/teamway-server/model"
)

// GetNextQuestion godoc
// @Summary      Gets next question
// @Description  Gets next unanswered question for a given test, or NoContent if there is none left
// @Tags         Take Test
// @Accept       json
// @Produce      json
// @Param        test-id     query     string     true  "Test Id"
// @Success      200  {string}  string "" Question
// @Success      204  {string}  string
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Failure      401      {object}  httputil.HTTPError
// @Failure      403      {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /test-taken/next [get]
func (c *Controller) GetNextQuestion(ctx *gin.Context) {
	var err error
	userId, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}
	testId := ctx.Query("test-id")

	id := model.NewTestTakenId(types.Id(testId), userId)

	test, err := model.TestOne(id.TestID)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}

	testTaken, err := model.TestTakenOneOrInsert(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}

	res := ""
	for _, question := range test.Questions {
		if _, ok := testTaken.Answers[question.Title]; !ok {
			res = question.Title
			break
		}
	}

	if res != "" {
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusNoContent, "")
	}
}

// AnswerQuestion godoc
// @Summary      Answer a question
// @Description  Saves an answer for a given question in a given test
// @Tags         Take Test
// @Accept       json
// @Produce      json
// @Param        test-id     query     string     true  "Test Id"
// @Param        question-title     query     string     true  "Question title"
// @Param        answer-title     query     string     true  "Answer title"
// @Success      200  {object}  model.TestTaken
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Failure      401      {object}  httputil.HTTPError
// @Failure      403      {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /test-taken [post]
func (c *Controller) AnswerQuestion(ctx *gin.Context) {
	var err error
	userId, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}
	testId := ctx.Query("test-id")
	questionTitle := ctx.Query("question-title")
	answerTitle := ctx.Query("answer-title")

	id := model.NewTestTakenId(types.Id(testId), userId)

	test, err := model.TestOne(id.TestID)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}

	testTaken, err := model.TestTakenOneOrInsert(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	testTakenCopy := testTaken.Copy()

	testTakenCopy.Answers[questionTitle] = answerTitle
	err = testTakenCopy.ValidationWithTest(test)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	err = model.TestTakenSave(testTakenCopy)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	res := *testTakenCopy

	ctx.JSON(http.StatusOK, res)
}

// GetTestResult godoc
// @Summary      Gets the test result
// @Description  Calculate and return the test result
// @Tags         Take Test
// @Accept       json
// @Produce      json
// @Param        test-id     query     string     true  "Test Id"
// @Success      200  {string}  string
// @Success      202  {string}  string
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Failure      403      {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /test-taken [get]
func (c *Controller) GetTestResult(ctx *gin.Context) {
	var err error
	userId, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}
	testId := ctx.Query("test-id")

	id := model.NewTestTakenId(types.Id(testId), userId)

	test, err := model.TestOne(id.TestID)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}

	testTaken, err := model.TestTakenOneOrInsert(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}

	err = model.TestTakenCalculateResult(testTaken, test)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	var res string
	if testTaken.Score != nil {
		if *testTaken.Score >= test.ScoreThreshold {
			res = test.PositiveResult
		} else {
			res = test.NegativeResult
		}
		ctx.JSON(http.StatusOK, res)
	} else {
		res = "Not all questions are answered"
		ctx.JSON(http.StatusAccepted, res)
	}
}
