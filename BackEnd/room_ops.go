package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomEvent struct {
	TimeStamp int64  `json:"time_stamp" bson:"time_stamp" yaml:"time_stamp"`
	Event     string `json:"event" bson:"event" yaml:"event"`
	Story     string `json:"event_story" bson:"event_story" yaml:"event_story"`
}

type Room struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id" yaml:"_id" `
	RoomId        string             `bson:"room_id" json:"room_id" yaml:"room_id"`
	Name          string             `json:"name" bson:"name" yaml:"name"`
	Private       bool               `json:"private" bson:"private" yaml:"private"`
	Level         string             `json:"level" bson:"level" yaml:"level"`
	Topics        []string           `json:"topics" bson:"topics" yaml:"topics"`
	Owner         string             `json:"owner" bson:"owner" yaml:"owner"`
	OwnerASMember bool               `json:"add_owner" bson:"add_owner" yaml:"add_owner"`
	Duration      int                `json:"duration" bson:"duration" yaml:"duration"`
	StartTime     int64              `json:"start_time" bson:"start_time" yaml:"start_time"`
	ProblemCount  int                `json:"problem_count" bson:"problem_count" yaml:"problem_count"`
	MaxMembers    int                `json:"max_members" bson:"max_members" yaml:"max_members"`
	TimeLine      []RoomEvent        `json:"timeline" bson:"timeline" yaml:"timeline"`
	Members       []string           `json:"members" bson:"members" yaml:"members"`
	Problems      []string           `json:"problems" bson:"problems" yaml:"problems"`
	Status        string             `json:"status" bson:"status" yaml:"status"`
	Results       map[string]int     `json:"results" bson:"results" yaml:"results"`
	CreatedAt     int64              `json:"created_at" bson:"created_at" yaml:"created_at"`
	UpdatedAt     int64              `json:"updated_at" bson:"updated_at" yaml:"updated_at"`
}

func create_room() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var new_room Room
		err := c.BindJSON(&new_room)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		current_time := time.Now().UnixNano()
		new_room.ID = primitive.NewObjectID()
		new_room.RoomId = new_room.ID.Hex()
		new_room.CreatedAt = current_time
		new_room.UpdatedAt = current_time
		new_room.Status = "Yet to Start"
		new_room.TimeLine = append(new_room.TimeLine, RoomEvent{Event: "Room Created", Story: new_room.Owner + " created this room", TimeStamp: current_time})
		if new_room.OwnerASMember {
			new_room.Members = append(new_room.Members, new_room.Owner)
		}
		room_collection := client.Database("BROCODE").Collection("rooms")
		insertionCount, err := room_collection.InsertOne(ctx, new_room)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			c.IndentedJSON(http.StatusCreated, insertionCount)
			return
		}
	}
}

func get_rooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		userID := c.Param("user_id")
		var to_return []Room
		room_collection := client.Database("BROCODE").Collection("rooms")
		result, err := room_collection.Find(ctx, bson.M{"$or": bson.M{"private": false, "owner": userID}})
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		err = result.All(ctx, &to_return)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"status": true, "data": to_return})
	}
}

func join_room() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		type Join_Room struct {
			UserId string `bson:"user_id" json:"user_id" yaml:"user_id"`
			RoomId string `bson:"room_id" json:"room_id" yaml:"room_id"`
		}
		var joining Join_Room
		err := c.BindJSON(&joining)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var room Room
		room_collection := client.Database("BROCODE").Collection("rooms")
		err = room_collection.FindOne(ctx, bson.M{"room_id": joining.RoomId}).Decode(&room)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		if room.MaxMembers > len(room.Members) {
			_, err := room_collection.UpdateOne(ctx, bson.M{"room_id": joining.RoomId}, bson.M{"$push": bson.M{"max_members": joining.UserId}})
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
				return
			}
			room.Members = append(room.Members, joining.UserId)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": true, "data": room})
			return
		}
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": "room is already full"})
	}
}

func leave_room() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		type Leave_Room struct {
			UserId string `bson:"user_id" json:"user_id" yaml:"user_id"`
			RoomId string `bson:"room_id" json:"room_id" yaml:"room_id"`
		}
		var leaving Leave_Room
		err := c.BindJSON(&leaving)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		room_collection := client.Database("BROCODE").Collection("rooms")
		_, err = room_collection.UpdateOne(ctx, bson.M{"room_id": leaving.RoomId}, bson.M{"$pull": bson.M{"max_members": leaving.UserId}})
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": true})
	}
}
