package load_quiz

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
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
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	mongourl := "mongodb://172.31.21.185:27017"
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongourl).SetAuth(credential))
	if err != nil {
		return nil, err
	}

	questionsCollection := client.Database("quiz").Collection("questions")

	var questions []Question

	if difficulty != "" {

		pipeline := []bson.M{
			{"$match": bson.M{
				"difficulty": difficulty}},
			{"$sample": bson.M{
				"size": 15}},
		}

		if topic != "" {
			pipeline = append(pipeline, bson.M{"$match": bson.M{"topic": topic}})
		}

		cursor, err := questionsCollection.Aggregate(ctx, pipeline)
		err = cursor.All(ctx, &questions)

		if err != nil {
			return nil, err
		}
	} else {

		var easyQuestions []Question
		var mediumQuestions []Question
		var hardQuestions []Question

		// Find easy difficulty
		easypipeline := []bson.M{
			{"$match": bson.M{
				"difficulty": "easy"}},
			{"$sample": bson.M{
				"size": 5}},
		}
		if topic != "" {
			easypipeline = append(easypipeline, bson.M{"$match": bson.M{"topic": topic}})
		}
		cursor, err := questionsCollection.Aggregate(ctx, easypipeline)
		err = cursor.All(ctx, &easyQuestions)
		if err != nil {
			return nil, err
		}

		// Find medium difficulty
		mediumpipeline := []bson.M{
			{"$match": bson.M{
				"difficulty": "medium"}},
			{"$sample": bson.M{
				"size": 5}},
		}
		if topic != "" {
			mediumpipeline = append(mediumpipeline, bson.M{"$match": bson.M{"topic": topic}})
		}
		cursor, err = questionsCollection.Aggregate(ctx, mediumpipeline)
		err = cursor.All(ctx, &mediumQuestions)
		if err != nil {
			return nil, err
		}

		// Find hard difficulty
		hardpipeline := []bson.M{
			{"$match": bson.M{
				"difficulty": "hard"}},
			{"$sample": bson.M{
				"size": 5}},
		}
		if topic != "" {
			hardpipeline = append(hardpipeline, bson.M{"$match": bson.M{"topic": topic}})
		}
		cursor, err = questionsCollection.Aggregate(ctx, hardpipeline)
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
