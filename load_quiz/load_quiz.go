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

func LoadQuizData() ([]Question, error) {

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
	filter := bson.D{{"topic", "arts"}}

	// retrieve all the documents that match the filter
	cursor, err := questionsCollection.Find(ctx, filter)

	var questions []Question
	err = cursor.All(ctx, &questions)

	if err != nil {
		return nil, err
	}

	return questions, nil
}
