package test

import (
	"github.com/oltur/teamway-server/model"
)

func getFirstTest() (test *model.Test, err error) {
	// get first test's Id
	tests, err := model.TestsAll()
	if err != nil {
		return
	}
	test = tests[0]
	return
}

func findQuestion(questions []*model.Question, title string) (res *model.Question) {
	for _, question := range questions {
		if question.Title == title {
			res = question
			return
		}
	}
	return
}
