package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	piston "github.com/milindmadhukar/go-piston"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RunBody struct {
	ID                 primitive.ObjectID `bson:"_id" json:"_id" yaml:"_id" `
	SubmissionID       string             `json:"submission_id" bson:"submission_id" yaml:"submission_id"`
	UserID             string             `json:"user_id" bson:"user_id" yaml:"user_id"`
	ProblemID          string             `json:"problem_id" bson:"problem_id" yaml:"problem_id"`
	Language           string             `json:"language" bson:"language" yaml:"language"`
	Version            string             `json:"version" bson:"version" yaml:"version"`
	FileName           string             `json:"file_name" bson:"file_name" yaml:"file_name"`
	Code               string             `json:"code" bson:"code" yaml:"code"`
	Stdin              string             `json:"stdin,omitempty" bson:"stdin" yaml:"stdin"`
	Args               []string           `json:"args,omitempty" bson:"args" yaml:"args"`
	CompileTimeout     int                `json:"compile_timeout,omitempty" bson:"compile_timeout" yaml:"compile_timeout"`
	RunTimeout         int                `json:"run_timeout,omitempty" bson:"run_timeout" yaml:"run_timeout"`
	CompileMemoryLimit int                `json:"compile_memory_limit,omitempty" bson:"compile_memory_limit" yaml:"compile_memory_limit"`
	RunMemoryLimit     int                `json:"run_memory_limit,omitempty" bson:"run_memory_limit" yaml:"run_memory_limit"`
	ResultCode         string             `json:"result_code" bson:"result_code" yaml:"result_code"`
}

func run_code() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var req_body RunBody
		err := c.BindJSON(&req_body)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		code := piston.Code{Content: req_body.Code}
		output, err := piston_client.Execute(req_body.Language, req_body.Version,
			[]piston.Code{code}, piston.Stdin(req_body.Stdin))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"status": true, "output": output.GetOutput()})
	}
}

func submit_code() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var req_body RunBody
		err := c.BindJSON(&req_body)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		submission_collection := client.Database("BROCODE").Collection("submissions")
		code := piston.Code{Content: req_body.Code}
		var problem Problems
		problem_collection := client.Database("BROCODE").Collection("problems")
		err = problem_collection.FindOne(ctx, bson.M{"problem_id": req_body.ProblemID}).Decode(&problem)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		var CA int
		type Failed_Case struct {
			TestCase string `json:"test_case"`
			Expected string `json:"expected"`
			Output   string `json:"output"`
		}
		var first_Failed Failed_Case
		for _, tc := range problem.TestCases {
			output, err := piston_client.Execute(req_body.Language, req_body.Version,
				[]piston.Code{code}, piston.Stdin(tc.Input))
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
				return
			}
			out_string := output.GetOutput()
			if tc.Answer == output.Run.Output {
				CA++
			} else {
				first_Failed = Failed_Case{TestCase: tc.Input, Expected: tc.Answer, Output: out_string}
			}
		}
		req_body.ID = primitive.NewObjectID()
		req_body.SubmissionID = req_body.ID.Hex()
		if CA == len(problem.TestCases) {
			req_body.ResultCode = "Correct Answer"
		} else {
			req_body.ResultCode = "Wrong Answer"
		}
		_, err = submission_collection.InsertOne(ctx, req_body)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		_, err = problem_collection.UpdateOne(ctx, bson.M{"problem_id": req_body.ProblemID}, bson.M{"$push": bson.M{"submissions": req_body.SubmissionID}})
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		if CA == len(problem.TestCases) {
			c.IndentedJSON(http.StatusOK, gin.H{"status": true, "code": "Correct Answer", "total_cases": len(problem.TestCases), "passed_cases": CA})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"status": true, "code": "Wrong Answer", "one_of_failed": first_Failed, "total_cases": len(problem.TestCases), "passed_cases": CA})
		}
	}
}
