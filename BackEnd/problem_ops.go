package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Submissions struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id" yaml:"_id" `
	SubmissionID string             `json:"submission_id" bson:"submission_id" yaml:"submission_id"`
	Language     string             `json:"lang" bson:"lang" yaml:"lang"`
	SubmittedBy  string             `json:"submitted_by" bson:"submitted_by" yaml:"submitted_by"`
	Code         string             `json:"code" bson:"code" yaml:"code"`
	CreatedAt    int64              `json:"created_at" bson:"created_at" yaml:"created_at"`
	ResultCode   string             `json:"res_code" bson:"res_code" yaml:"res_code"`
	TimeTaken    int64              `json:"time" bson:"time" yaml:"time"`
	SpaceUsed    float32            `json:"space" bson:"space" yaml:"space"`
}

type TestCase struct {
	Input  string `json:"input" bson:"input" yaml:"input"`
	Answer string `json:"answer" bson:"answer" yaml:"answer"`
}

type Problems struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id" yaml:"_id" `
	ProblemID   string             `json:"problem_id" bson:"problem_id" yaml:"problem_id"`
	Level       string             `json:"level" bson:"level" yaml:"level"`
	Topics      []string           `json:"topics" bson:"topics" yaml:"topics"`
	Title       string             `json:"title" bson:"title" yaml:"title"`
	CreatedBy   string             `json:"created_by" bson:"created_by" yaml:"created_by"`
	Description string             `json:"p_description" bson:"p_description" yaml:"p_description"`
	TestCases   []TestCase         `json:"test_cases" bson:"test_cases" yaml:"test_cases"`
	Submissions []string           `json:"submissions" bson:"submissions" yaml:"submissions"`
	CreatedAt   int64              `json:"created_at" bson:"created_at" yaml:"created_at"`
	UpdatedAt   int64              `json:"updated_at" bson:"updated_at" yaml:"updated_at"`
	Verified    bool               `json:"verified" bson:"verified" yaml:"verified"`
}

// func test() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	problem_collection := client.Database("BROCODE").Collection("problems")
// 	update := bson.A{
// 		bson.M{"$set": bson.M{
// 			"verified": bson.M{
// 				"$cond": bson.A{
// 					bson.D{{Key: "$eq", Value: bson.A{"$verified", true}}}, false, true,
// 				},
// 			},
// 		},
// 		},
// 	}
// 	_, err := problem_collection.UpdateOne(ctx, bson.M{"title": "Jump game"}, update)
// 	if err != nil {

// 	}

// }

func create_problem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var new_problem Problems
		err := c.BindJSON(&new_problem)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		current_time := time.Now().UnixNano()
		new_problem.ID = primitive.NewObjectID()
		new_problem.ProblemID = new_problem.ID.Hex()
		new_problem.CreatedAt = current_time
		new_problem.UpdatedAt = current_time

		var user User
		user_collection := client.Database("BROCODE").Collection("users")
		if err = user_collection.FindOne(ctx, bson.M{"name": new_problem.CreatedBy}).Decode(&user); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if user.Role == "ADMIN" || user.Role == "PRO" {
			new_problem.Verified = true
		} else {
			new_problem.Verified = false
		}

		problem_collection := client.Database("BROCODE").Collection("problems")
		insertionCount, err := problem_collection.InsertOne(ctx, new_problem)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			c.IndentedJSON(http.StatusCreated, insertionCount)
			return
		}
	}
}

func get_problems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		type Filter_Problems struct {
			Level  string   `json:"level" bson:"level" yaml:"level"`
			Topics []string `json:"topics" bson:"topics" yaml:"topics"`
			Title  string   `json:"title" bson:"title" yaml:"title"`
		}
		var to_return []Problems
		var filter Filter_Problems
		err := c.BindJSON(&filter)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
			return
		}
		problem_collection := client.Database("BROCODE").Collection("problems")
		result, err := problem_collection.Find(ctx, bson.M{"verified": true})
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
			return
		}
		err = result.All(ctx, &to_return)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"status": true, "data": to_return})
	}
}
