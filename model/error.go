package model

import "errors"

type ErrTypeNotFound error
type ErrTypeInvalidID error
type ErrTypeInvalidUserName error
type ErrTypeInvalidPassword error
type ErrTypeUserNotFoundInContext error
type ErrTypeCannotGenerateUserToken error
type ErrTypeCannotValidateUserToken error
type ErrTypeAccessDenied error
type ErrTypeUnauthorized error
type ErrTypeIdExists error
type ErrTypeUserNameExists error
type ErrTypeActiveSessionExists error
type ErrTypeInvalidTestTitle error
type ErrTypeInvalidQuestionTitle error
type ErrTypeInvalidAnswerTitle error
type ErrTypeInvalidPositiveResult error
type ErrTypeInvalidNegativeResult error
type ErrTypeInvalidScoreThreshold error
type ErrTypeInvalidQuestions error
type ErrTypeInvalidAnswers error
type ErrTypeInvalidScore error
type ErrTypeInvalidUserId error
type ErrTypeInvalidTestId error
type ErrTypeNoTestForValidation error
type ErrTypeAnswersDoNotMatchTest error
type ErrTypeDuplicateQuestion error
type ErrTypeDuplicateAnswer error

var (
	ErrNotFound                = ErrTypeNotFound(errors.New("not found"))
	ErrInvalidID               = ErrTypeInvalidID(errors.New("invalid id"))
	ErrInvalidUserName         = ErrTypeInvalidUserName(errors.New("invalid user name"))
	ErrInvalidPassword         = ErrTypeInvalidPassword(errors.New("invalid password"))
	ErrUserNotFoundInContext   = ErrTypeUserNotFoundInContext(errors.New("user not found in context"))
	ErrCannotGenerateUserToken = ErrTypeCannotGenerateUserToken(errors.New("cannot generate user token"))
	ErrCannotValidateUserToken = ErrTypeCannotValidateUserToken(errors.New("cannot validate user token"))
	ErrAccessDenied            = ErrTypeAccessDenied(errors.New("access denied"))
	ErrUnauthorized            = ErrTypeUnauthorized(errors.New("'Authorization' is required Header"))
	ErrIdExists                = ErrTypeIdExists(errors.New("this ID already exists"))
	ErrUserNameExists          = ErrTypeUserNameExists(errors.New("user with given name already exists"))
	ErrActiveSessionExists     = ErrTypeActiveSessionExists(errors.New("there is already an active session using your account. Use /user/logout/all API to drop it"))
	ErrInvalidTestTitle        = ErrTypeInvalidTestTitle(errors.New("invalid test title"))
	ErrInvalidQuestionTitle    = ErrTypeInvalidQuestionTitle(errors.New("invalid question title"))
	ErrInvalidAnswerTitle      = ErrTypeInvalidAnswerTitle(errors.New("invalid answer title"))
	ErrInvalidPositiveResult   = ErrTypeInvalidPositiveResult(errors.New("invalid positive result"))
	ErrInvalidNegativeResult   = ErrTypeInvalidNegativeResult(errors.New("invalid negative result"))
	ErrInvalidScoreThreshold   = ErrTypeInvalidScoreThreshold(errors.New("invalid score threshold"))
	ErrInvalidQuestions        = ErrTypeInvalidQuestions(errors.New("invalid questions"))
	ErrInvalidAnswers          = ErrTypeInvalidAnswers(errors.New("invalid answers"))
	ErrInvalidScore            = ErrTypeInvalidScore(errors.New("invalid answer score"))
	ErrInvalidUserId           = ErrTypeInvalidUserId(errors.New("invalid userId"))
	ErrInvalidTestId           = ErrTypeInvalidTestId(errors.New("invalid testId"))
	ErrNoTestForValidation     = ErrTypeNoTestForValidation(errors.New("no test is provided for validation of test result"))
	ErrAnswersDoNotMatchTest   = ErrTypeAnswersDoNotMatchTest(errors.New("there are too many answers or answered questions or answers do not match test questions"))
	ErrDuplicateQuestion       = ErrTypeDuplicateQuestion(errors.New("duplicate question"))
	ErrDuplicateAnswer         = ErrTypeDuplicateAnswer(errors.New("duplicate answer"))
)
