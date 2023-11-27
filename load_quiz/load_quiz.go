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

		sampleStage := bson.D{{"$sample", bson.D{{"size", 15}}}}

		pipeline := mongo.Pipeline{}
		if topic != "" {
			pipeline = append(pipeline, bson.D{{"$match", bson.M{"topic": topic}}})
		}
		pipeline = append(pipeline, bson.D{{"$match", bson.D{{"difficulty", difficulty}}}})
		pipeline = append(pipeline, sampleStage)

		opts := options.Aggregate().SetMaxTime(2 * time.Second)
		cursor, err := questionsCollection.Aggregate(ctx, pipeline, opts)
		if err != nil {
			return nil, err
		}

		err = cursor.All(ctx, &questions)
		if err != nil {
			return nil, err
		}
	} else {

		var easyQuestions []Question
		var mediumQuestions []Question
		var hardQuestions []Question

		sampleStage := bson.D{{"$sample", bson.D{{"size", 5}}}}

		easyPipeline := mongo.Pipeline{}
		if topic != "" {
			easyPipeline = append(easyPipeline, bson.D{{"$match", bson.M{"topic": topic}}})
		}
		easyPipeline = append(easyPipeline, bson.D{{"$match", bson.D{{"difficulty", "easy"}}}})
		easyPipeline = append(easyPipeline, sampleStage)

		opts := options.Aggregate().SetMaxTime(2 * time.Second)
		cursor, err := questionsCollection.Aggregate(ctx, easyPipeline, opts)
		if err != nil {
			return nil, err
		}

		err = cursor.All(ctx, &easyQuestions)
		if err != nil {
			return nil, err
		}

		mediumPipeline := mongo.Pipeline{}
		if topic != "" {
			easyPipeline = append(mediumPipeline, bson.D{{"$match", bson.M{"topic": topic}}})
		}
		easyPipeline = append(mediumPipeline, bson.D{{"$match", bson.D{{"difficulty", "medium"}}}})
		easyPipeline = append(mediumPipeline, sampleStage)

		cursor, err = questionsCollection.Aggregate(ctx, mediumPipeline, opts)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		err = cursor.All(ctx, &mediumQuestions)
		if err != nil {
			return nil, err
		}

		hardPipeline := mongo.Pipeline{}
		if topic != "" {
			easyPipeline = append(mediumPipeline, bson.D{{"$match", bson.M{"topic": topic}}})
		}
		easyPipeline = append(mediumPipeline, bson.D{{"$match", bson.D{{"difficulty", "hard"}}}})
		easyPipeline = append(mediumPipeline, sampleStage)

		cursor, err = questionsCollection.Aggregate(ctx, hardPipeline, opts)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
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
