package load_quiz

type Question struct {
	Id         int32           `bson:"id" json:"id"`
	En         QuestionDetails `bson:"en" json:"en"`
	De         QuestionDetails `bson:"de" json:"de"`
	Hu         QuestionDetails `bson:"hu" json:"hu"`
	Difficulty string          `bson:"difficulty" json:"difficulty"`
	Topic      string          `bson:"topic" json:"topic"`
}

type QuestionDetails struct {
	Text               string   `bson:"text" json:"text"`
	Answers            []string `bson:"answers" json:"answers"`
	CorrectAnswerIndex int32    `bson:"correct_answer_index" json:"correct_answer_index"`
}

type RequestBody struct {
	EasyQuestionExceptions   []int32 `bson:"easyQuestions" json:"easyQuestions"`
	MediumQuestionExceptions []int32 `bson:"mediumQuestions" json:"mediumQuestions"`
	HardQuestionExceptions   []int32 `bson:"hardQuestions" json:"hardQuestions"`
}

type LoadResponse struct {
	Questions []Question `bson:"questions" json:"questions"`
	Exception Exception  `bson:"exception" json:"exception"`
}

type Exception struct {
	ResetEasyFilter   bool `bson:"resetEasyFilter" json:"resetEasyFilter"`
	ResetMediumFilter bool `bson:"resetMediumFilter" json:"resetMediumFilter"`
	ResetHardFilter   bool `bson:"resetHardFilter" json:"resetHardFilter"`
}
