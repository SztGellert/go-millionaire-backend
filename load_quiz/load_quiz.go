package load_quiz

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

func ConnectMongo() {

	// test db with private network access
	credential := options.Credential{
		Username: "admin",
		Password: "admin",
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://172.31.21.185:27017").SetAuth(credential))
	if err != nil {
		fmt.Println(err.Error())
		log.Println(err)
		panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println(err.Error())
		log.Println(err)
		panic(err)
	}
}

func LoadQuizData(topic string, difficulty string) ([]Question, error) {

	credential := options.Credential{
		Username: "admin",
		Password: "admin",
	}

	ctx := context.TODO()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://172.31.21.185:27017").SetAuth(credential))
	if err != nil {
		return nil, err
	}

	questionsCollection := client.Database("quiz").Collection("questions")
	var questions []Question

	if difficulty != "" {
		opts := options.Find().SetLimit(15)

		filter := bson.D{{"topic", topic}, {"difficulty", difficulty}}
		// retrieve all the documents that match the filter
		cursor, err := questionsCollection.Find(ctx, filter, opts)
		err = cursor.All(ctx, &questions)

		if err != nil {
			return nil, err
		}
	} else {

		opts := options.Find().SetLimit(5)

		// Find easy difficulty
		var easyQuestions []Question
		filter := bson.D{{"topic", topic}, {"difficulty", "easy"}}
		cursor, err := questionsCollection.Find(ctx, filter, opts)
		err = cursor.All(ctx, &easyQuestions)

		// Find medium difficulty
		var mediumQuestions []Question
		filter = bson.D{{"topic", topic}, {"difficulty", "medium"}}
		cursor, err = questionsCollection.Find(ctx, filter, opts)
		err = cursor.All(ctx, &mediumQuestions)

		// Find hard difficulty
		var hardQuestions []Question
		filter = bson.D{{"topic", topic}, {"difficulty", "hard"}}
		cursor, err = questionsCollection.Find(ctx, filter, opts)
		err = cursor.All(ctx, &hardQuestions)

		if err != nil {
			return nil, err
		}

		questions = append(questions, easyQuestions...)
		questions = append(questions, mediumQuestions...)
		questions = append(questions, hardQuestions...)

	}

	return questions, nil
}
