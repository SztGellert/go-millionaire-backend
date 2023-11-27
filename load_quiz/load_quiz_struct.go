package load_quiz

type Question struct {
	Id            int32    `bson:"id" json:"id"`
	Value         string   `bson:"value" json:"value"`
	Answers       []string `bson:"answers" json:"answers"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`
	Difficulty    string   `bson:"difficulty" json:"difficulty"`
	Topic         string   `bson:"topic" json:"topic"`
}
