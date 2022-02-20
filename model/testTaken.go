package model

import (
	"errors"
	"github.com/oltur/teamway-server/types"
)

type TestTakenId struct {
	TestID types.Id `json:"testId"`
	UserID types.Id `json:"userId"`
}

func NewTestTakenId(testId types.Id, userId types.Id) (res TestTakenId) {
	return TestTakenId{
		TestID: testId,
		UserID: userId,
	}
}

type TestTaken struct {
	ID      TestTakenId       `json:"testTakenId"`
	Answers map[string]string `json:"answers"` // key is question title, value is answer title
	Score   *int              `json:"score"`
}

func (t TestTaken) Copy() (t2 *TestTaken) {
	answers := make(map[string]string, len(t.Answers))
	for k, v := range t.Answers {
		answers[k] = v
	}
	var score *int
	if t.Score != nil {
		x := *t.Score
		score = &x
	}
	t2 = &TestTaken{
		ID:      t.ID,
		Answers: answers,
		Score:   score,
	}
	return
}

func (t TestTaken) Validation() (err error) {
	if t.ID.UserID == "" {
		err = ErrInvalidUserId
		return
	}
	if t.ID.TestID == "" {
		err = ErrInvalidTestId
		return
	}
	for k, v := range t.Answers {
		if k == "" {
			err = ErrInvalidQuestionTitle
			return
		}
		if v == "" {
			err = ErrInvalidAnswerTitle
			return
		}
	}
	return
}

func (t TestTaken) ValidationWithTest(test *Test) (err error) {
	err = t.Validation()
	if err != nil {
		return
	}
	if test == nil {
		err = ErrNoTestForValidation
		return
	}

	if len(t.Answers) > len(test.Questions) {
		err = ErrAnswersDoNotMatchTest
		return
	}

	questionsAndAnswersMap := make(map[string]map[string]interface{}, len(test.Questions))
	for _, question := range test.Questions {
		questionsAndAnswersMap[question.Title] = make(map[string]interface{}, len(question.Answers))
		for _, answer := range question.Answers {
			questionsAndAnswersMap[question.Title][answer.Title] = nil
		}
	}
	for q, a := range t.Answers {
		question, ok1 := questionsAndAnswersMap[q]
		if !ok1 {
			err = ErrAnswersDoNotMatchTest
			return
		}
		_, ok2 := question[a]
		if !ok2 {
			err = ErrAnswersDoNotMatchTest
			return
		}
	}
	return
}

func TestTakenCalculateResult(testTaken *TestTaken, test *Test) (err error) {
	if testTaken == nil || test == nil {
		err = ErrNotFound
		return
	}

	// do validation to avoid further error checking
	err = testTaken.ValidationWithTest(test)
	if err != nil {
		return
	}

	questionsAndAnswersMap := make(map[string]map[string]int, len(test.Questions))
	for _, question := range test.Questions {
		questionsAndAnswersMap[question.Title] = make(map[string]int, len(question.Answers))
		for _, answer := range question.Answers {
			questionsAndAnswersMap[question.Title][answer.Title] = answer.Score
		}
	}

	// not fully answered
	if len(testTaken.Answers) != len(test.Questions) {
		testTaken.Score = nil
		return
	}

	totalScore := 0
	for q, a := range testTaken.Answers {
		question, ok := questionsAndAnswersMap[q]
		if !ok {
			break
		}
		score, ok2 := question[a]
		if !ok2 {
			break
		}
		totalScore = totalScore + score
	}
	testTaken.Score = &totalScore

	err = TestTakenSave(testTaken)
	if err != nil {
		return
	}

	return
}

func TestTakenOne(id TestTakenId) (res *TestTaken, err error) {
	res, ok := testsTakenByIds[id]
	if ok {
		return
	}
	err = ErrNotFound
	return
}

func TestTakenOneOrInsert(id TestTakenId) (res *TestTaken, err error) {
	res, err = TestTakenOne(id)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return
		}
	} else {
		return
	}
	res = &TestTaken{
		ID:      id,
		Answers: make(map[string]string),
		Score:   nil,
	}
	res, err = TestTakenInsert(res)
	if err != nil {
		res = nil
		return
	}
	return
}

// TestTakenSave Internal use only
func TestTakenSave(req *TestTaken) (err error) {
	testsTakenByIds[req.ID] = req
	return
}

//func TestTakenDelete(id TestTakenId) (err error) {
//	if _, ok := testsTakenByIds[id]; ok {
//		err = ErrNotFound
//		return
//	}
//	delete(testsTakenByIds, id)
//	return
//}

var testsTakenByIds = make(map[TestTakenId]*TestTaken)

//func GetMapKeysForTestsTaken(m map[types.Id]*TestTaken) (res []types.Id) {
//	res = make([]types.Id, len(m))
//	i := 0
//	for k := range m {
//		res[i] = k
//		i++
//	}
//	return
//}

//func GetMapValuesForTestsTaken(m map[types.Id]*TestTaken) (res []*TestTaken) {
//	res = make([]*TestTaken, len(m))
//	i := 0
//	for _, v := range m {
//		res[i] = v
//		i++
//	}
//	return res
//}

func TestTakenInsert(req *TestTaken) (res *TestTaken, err error) {
	_, err = TestTakenOne(req.ID)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return
		} else {
			err = nil
		}
	} else {
		err = ErrIdExists
		return
	}

	testsTakenByIds[req.ID] = req
	res = req
	return
}
