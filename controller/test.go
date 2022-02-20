package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/oltur/teamway-server/httputil"
	"github.com/oltur/teamway-server/model"
	"github.com/oltur/teamway-server/types"
	"net/http"
)

// ShowTest godoc
// @Summary      Show a test
// @Description  get string by ID
// @Tags         Test
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Test ID"
// @Success      200  {object}  model.Test
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /test/{id} [get]
func (c *Controller) ShowTest(ctx *gin.Context) {
	s := ctx.Param("id")
	id := types.Id(s)
	test, err := model.TestOne(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, test)
}

// ListTests godoc
// @Summary      List tests
// @Description  get tests
// @Tags         Test
// @Accept       json
// @Produce      json
// @Success      200  {array}   model.Test
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /test [get]
func (c *Controller) ListTests(ctx *gin.Context) {
	tests, err := model.TestsAll()
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, tests)
}

// AddTest godoc
// @Summary      Add test
// @Description  Add new test
// @Tags         Test
// @Accept       json
// @Produce      json
// @Param        test  body      model.AddTestRequest  true  "Add test request"
// @Success      200      {object}  model.Test
// @Failure      400      {object}  httputil.HTTPError
// @Failure      404      {object}  httputil.HTTPError
// @Failure      500      {object}  httputil.HTTPError
// @Failure      401      {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /test [post]
func (c *Controller) AddTest(ctx *gin.Context) {
	_, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}

	var req model.AddTestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := model.TestInsert(&req)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// DeleteTest godoc
// @Summary      Delete a test
// @Description  Delete by test ID
// @Tags         Test
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Test ID"
// @Success 	 204  {string} string "Ok"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Failure      401      {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /test/{id} [delete]
func (c *Controller) DeleteTest(ctx *gin.Context) {
	s := ctx.Param("id")
	id := types.Id(s)

	// delete
	err := model.TestDelete(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusNoContent, "Ok")
}
