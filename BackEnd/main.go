package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	piston "github.com/milindmadhukar/go-piston"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var piston_client *piston.Client

func mongo_initializer() *mongo.Client {
	m_client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = m_client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("unable to ping", err)
	}
	log.Println("Mongo client created")
	return m_client
}

func piston_initializer() *piston.Client {
	return piston.CreateDefaultClient()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	client = mongo_initializer()
	piston_client = piston_initializer()
	// test()

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/dbs", show_dbs())

	router.POST("/signup", create_user())
	router.POST("/login", login())

	router.POST("/user/:user_id", get_user())

	router.POST("/createRoom", create_room())
	router.POST("/room/:user_id", get_rooms())
	router.POST("/room/join", join_room())
	router.POST("/room/leave", leave_room())

	router.POST("/createProblem", create_problem())
	router.POST("/get_problems", get_problems())

	router.POST("/run", run_code())
	router.POST("/submit", submit_code())

	router.Run("localhost:8080")
}
