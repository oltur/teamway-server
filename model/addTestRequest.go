package model

type AddTestRequest struct {
	Title          string      `json:"title"`
	Questions      []*Question `json:"questions"`
	ScoreThreshold int         `json:"scoreThreshold"`
	PositiveResult string      `json:"positiveResult"`
	NegativeResult string      `json:"negativeResult"`
}

func (t AddTestRequest) Validation() (err error) {
	if t.Title == "" {
		err = ErrInvalidTestTitle
		return
	}
	if t.PositiveResult == "" {
		err = ErrInvalidPositiveResult
		return
	}
	if t.NegativeResult == "" {
		err = ErrInvalidNegativeResult
		return
	}
	if t.ScoreThreshold <= 0 {
		err = ErrInvalidScoreThreshold
		return
	}
	if len(t.Questions) == 0 {
		err = ErrInvalidQuestions
		return
	}
	questionsAndAnswersMap := make(map[string]map[string]interface{}, len(t.Questions))
	for _, q := range t.Questions {
		if q.Title == "" {
			err = ErrInvalidQuestionTitle
			return
		}
		if len(q.Answers) == 0 {
			err = ErrInvalidAnswers
			return
		}
		if _, ok := questionsAndAnswersMap[q.Title]; ok {
			err = ErrDuplicateQuestion
			return
		}
		questionsAndAnswersMap[q.Title] = make(map[string]interface{})
		for _, a := range q.Answers {
			if a.Title == "" {
				err = ErrInvalidAnswerTitle
				return
			}
			if a.Score <= 0 {
				err = ErrInvalidScore
				return
			}
			if _, ok := questionsAndAnswersMap[q.Title][a.Title]; ok {
				err = ErrDuplicateAnswer
				return
			}
			questionsAndAnswersMap[q.Title][a.Title] = nil
		}
	}
	return
}
