package model

type GetNextQuestionResponse struct {
	TestFinished   bool   `json:"testFinished"`
	Question       string `json:"question"`
	TotalQuestions int    `json:"totalQuestions"`
	QuestionNumber int    `json:"questionNumber"`
}
