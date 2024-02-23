package load_quiz

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

func ConnectMongo() {

	mongoUser := os.Getenv("MONGODB_USER")
	mongoPassword := os.Getenv("MONGODB_PASSWORD")

	// test db with private network access
	credential := options.Credential{
		Username: mongoUser,
		Password: mongoPassword,
	}

	mongoURI := os.Getenv("MONGODB_URI")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI).SetAuth(credential))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
}

func LoadQuizData(topic string, difficulty string, exceptions string) (LoadResponse, error) {

	mongoUser := os.Getenv("MONGODB_USER")
	mongoPassword := os.Getenv("MONGODB_PASSWORD")

	// test db with private network access
	credential := options.Credential{
		Username: mongoUser,
		Password: mongoPassword,
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	mongoURI := os.Getenv("MONGODB_URI")
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI).SetAuth(credential))
	if err != nil {
		return LoadResponse{}, err
	}
	exceptionRequest := RequestBody{}
	if exceptions != "" {
		err = json.Unmarshal([]byte(exceptions), &exceptionRequest)
		if err != nil {
			return LoadResponse{}, err
		}
	}

	mongoDatabase := os.Getenv("MONGODB_DATABASE")
	mongoCollection := os.Getenv("MONGODB_COLLECTION")

	questionsCollection := client.Database(mongoDatabase).Collection(mongoCollection)

	var questions []Question
	var resetEasyFilter, resetMediumFilter, resetHardFilter bool

	if difficulty != "" {
		var exceptionsList []int32
		var resetFilter bool
		if exceptions != "" {
			switch difficulty {
			case "easy":
				exceptionsList = exceptionRequest.EasyQuestionExceptions
				resetEasyFilter = resetFilter
			case "medium":
				exceptionsList = exceptionRequest.MediumQuestionExceptions
				resetMediumFilter = resetFilter
			default:
				exceptionsList = exceptionRequest.HardQuestionExceptions
				resetHardFilter = resetFilter
			}
		}

		sampleStage := bson.D{{"$sample", bson.D{{"size", 15}}}}

		pipeline := mongo.Pipeline{}
		var filter bson.D
		var countOpts *options.CountOptions
		var count int64

		if exceptions != "" && exceptionsList != nil {
			filter = bson.D{{"id", bson.D{{"$nin", exceptionsList}}}}
			count, err = questionsCollection.CountDocuments(context.TODO(), filter, countOpts)
			if err != nil {
				return LoadResponse{}, err
			}
			if count >= 15 {
				pipeline = append(pipeline, bson.D{{"$match", filter}})

			} else {
				resetFilter = true
			}
		}

		if topic != "" {
			pipeline = append(pipeline, bson.D{{"$match", bson.M{"topic": topic}}})
		}
		pipeline = append(pipeline, bson.D{{"$match", bson.D{{"difficulty", difficulty}}}})
		pipeline = append(pipeline, sampleStage)

		opts := options.Aggregate().SetMaxTime(2 * time.Second)
		cursor, err := questionsCollection.Aggregate(ctx, pipeline, opts)
		if err != nil {
			return LoadResponse{}, err
		}

		err = cursor.All(ctx, &questions)
		if err != nil {
			return LoadResponse{}, err
		}
	} else {

		var easyQuestions []Question
		var mediumQuestions []Question
		var hardQuestions []Question

		sampleStage := bson.D{{"$sample", bson.D{{"size", 5}}}}

		easyPipeline := mongo.Pipeline{}
		var filter bson.D
		var countOpts *options.CountOptions
		var count int64

		if exceptions != "" && exceptionRequest.EasyQuestionExceptions != nil {
			filter = bson.D{{"id", bson.D{{"$nin", exceptionRequest.EasyQuestionExceptions}}}}
			countOpts = options.Count().SetMaxTime(2 * time.Second)
			count, err = questionsCollection.CountDocuments(context.TODO(), filter, countOpts)
			if err != nil {
				return LoadResponse{}, err
			}
			if count >= 15 {
				easyPipeline = append(easyPipeline, bson.D{{"$match", filter}})
			} else {
				resetEasyFilter = true
			}
		}
		if topic != "" {
			easyPipeline = append(easyPipeline, bson.D{{"$match", bson.M{"topic": topic}}})
		}
		easyPipeline = append(easyPipeline, bson.D{{"$match", bson.D{{"difficulty", "easy"}}}})
		easyPipeline = append(easyPipeline, sampleStage)

		opts := options.Aggregate().SetMaxTime(2 * time.Second)
		cursor, err := questionsCollection.Aggregate(ctx, easyPipeline, opts)
		if err != nil {
			return LoadResponse{}, err
		}

		err = cursor.All(ctx, &easyQuestions)
		if err != nil {
			return LoadResponse{}, err
		}

		mediumPipeline := mongo.Pipeline{}
		if exceptions != "" && exceptionRequest.MediumQuestionExceptions != nil {
			filter = bson.D{{"id", bson.D{{"$nin", exceptionRequest.MediumQuestionExceptions}}}}
			count, err = questionsCollection.CountDocuments(context.TODO(), filter, countOpts)
			if err != nil {
				return LoadResponse{}, err
			}
			if count >= 15 {
				mediumPipeline = append(mediumPipeline, bson.D{{"$match", filter}})
			} else {
				resetMediumFilter = true
			}
		}
		if topic != "" {
			mediumPipeline = append(mediumPipeline, bson.D{{"$match", bson.M{"topic": topic}}})
		}
		mediumPipeline = append(mediumPipeline, bson.D{{"$match", bson.D{{"difficulty", "medium"}}}})
		mediumPipeline = append(mediumPipeline, sampleStage)

		cursor, err = questionsCollection.Aggregate(ctx, mediumPipeline, opts)
		if err != nil {
			return LoadResponse{}, err
		}
		err = cursor.All(ctx, &mediumQuestions)
		if err != nil {
			return LoadResponse{}, err
		}

		hardPipeline := mongo.Pipeline{}
		if exceptions != "" && exceptionRequest.HardQuestionExceptions != nil {
			filter = bson.D{{"id", bson.D{{"$nin", exceptionRequest.HardQuestionExceptions}}}}
			count, err = questionsCollection.CountDocuments(context.TODO(), filter, countOpts)
			if err != nil {
				return LoadResponse{}, err
			}
			if count >= 15 {
				hardPipeline = append(hardPipeline, bson.D{{"$match", filter}})
			} else {
				resetHardFilter = true
			}
		}
		if topic != "" {
			hardPipeline = append(hardPipeline, bson.D{{"$match", bson.M{"topic": topic}}})
		}
		hardPipeline = append(hardPipeline, bson.D{{"$match", bson.D{{"difficulty", "hard"}}}})
		hardPipeline = append(hardPipeline, sampleStage)

		cursor, err = questionsCollection.Aggregate(ctx, hardPipeline, opts)
		if err != nil {
			return LoadResponse{}, err
		}
		err = cursor.All(ctx, &hardQuestions)
		if err != nil {
			return LoadResponse{}, err
		}

		questions = append(questions, easyQuestions...)
		questions = append(questions, mediumQuestions...)
		questions = append(questions, hardQuestions...)

	}

	return LoadResponse{Questions: questions, Exception: Exception{ResetEasyFilter: resetEasyFilter, ResetMediumFilter: resetMediumFilter, ResetHardFilter: resetHardFilter}}, nil
}
