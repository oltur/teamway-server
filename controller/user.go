package controller

import (
	"errors"
	"github.com/oltur/teamway-server/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oltur/teamway-server/httputil"
	"github.com/oltur/teamway-server/model"
	"github.com/rs/xid"
)

// Login godoc
// @Summary      Login
// @Description  Logs user in
// @Tags         User
// @Accept       json
// @Produce      json
// @Param		 credentials body	model.LoginRequest true  "Login Request"
// @Success      200  {string}  model.LoginResponse
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      409  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /user/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	var err error
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	gwtToken, tokenExpires, err := c.DoLogin(req.UserName, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			httputil.NewError(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, model.ErrActiveSessionExists) {
			httputil.NewError(ctx, http.StatusConflict, err)
		} else {
			httputil.NewError(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	res := &model.LoginResponse{
		Token:        gwtToken,
		TokenExpires: tokenExpires,
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *Controller) DoLogin(userName string, password string) (gwtToken string, tokenExpires int64, err error) {
	user, err := model.GetUserByCredentials(userName, password)
	if err != nil {
		err = model.ErrNotFound
		return "", 0, err
	}

	token := xid.New().String()
	tokenExpires = time.Now().Add(30 * 24 * time.Hour).UnixMilli()

	gwtToken, err = c.createGwt(string(user.ID), token, tokenExpires)
	if err != nil {
		err = model.ErrCannotGenerateUserToken
		return "", 0, err
	}

	now := time.Now().UnixMilli()
	if user.Token != "" && user.TokenExpires >= now {
		err = model.ErrActiveSessionExists
		return "", 0, err
	}

	user.Token = token
	user.TokenExpires = tokenExpires

	err = model.UserSave(user)
	if err != nil {
		return "", 0, err
	}
	return gwtToken, tokenExpires, err
}

// Logout godoc
// @Summary      Logout
// @Description  Logs user out
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      204  {string} string "Ok"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Failure      401      {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /user/logout [post]
func (c *Controller) Logout(ctx *gin.Context) {
	userId, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}
	// check Admin role
	currentUser, err := model.UserOne(userId)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	err = model.UserLogout(currentUser.ID)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusNoContent, "Ok")
}

// LogoutAll godoc
// @Summary      Log out all user's sessions
// @Description  Logs current user ouy of all sessions
// @Tags         User
// @Accept       json
// @Produce      json
// @Param		 credentials body	model.LoginRequest true  "Login Request"
// @Success      204  {string}  string "Ok"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      409  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /user/logout/all [post]
func (c *Controller) LogoutAll(ctx *gin.Context) {
	var err error
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	user, err := model.GetUserByCredentials(req.UserName, req.Password)
	if err != nil {
		err = model.ErrNotFound
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}

	user.Token = ""
	user.TokenExpires = 0

	err = model.UserSave(user)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusNotFound, "Ok")
}

// ShowUser godoc
// @Summary      Show a user
// @Description  get string by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  model.User
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Failure      401      {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /user/{id} [get]
func (c *Controller) ShowUser(ctx *gin.Context) {
	s := ctx.Param("id")
	id := types.Id(s)

	userId, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}
	// can be viewed by themselves or by admin
	if userId != id {
		err = model.ErrAccessDenied
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}

	user, err := model.UserOne(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// AddUser godoc
// @Summary      Register
// @Description  Registers a user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      model.AddUserRequest  true  "Add user request"
// @Success      200      {object}  model.User
// @Failure      400      {object}  httputil.HTTPError
// @Failure      404      {object}  httputil.HTTPError
// @Failure      500      {object}  httputil.HTTPError
// @Router       /user [post]
func (c *Controller) AddUser(ctx *gin.Context) {
	var req model.AddUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := model.UserInsert(&req)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Update by json user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      model.UpdateUserRequest  true  "Update user info"
// @Success      200      {object}  model.UpdateUserRequest
// @Failure      400      {object}  httputil.HTTPError
// @Failure      404      {object}  httputil.HTTPError
// @Failure      500      {object}  httputil.HTTPError
// @Failure      401      {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /user/{id} [patch]
func (c *Controller) UpdateUser(ctx *gin.Context) {
	var updateUserRequest model.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&updateUserRequest); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	userId, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}
	// can be updated by themselves or by admin
	if userId != updateUserRequest.ID {
		err = model.ErrAccessDenied
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}

	err = model.UserUpdate(&updateUserRequest)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, updateUserRequest)
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete by user ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success 	 204  {string} string "Ok"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Failure      401      {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /user/{id} [delete]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	s := ctx.Param("id")
	id := types.Id(s)

	userId, err := c.getUserIdFromContext(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}
	// can be deleted by themselves or by admin
	if userId != id {
		err = model.ErrAccessDenied
		httputil.NewError(ctx, http.StatusForbidden, err)
		return
	}

	err = model.UserDelete(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusNoContent, "Ok")
}
