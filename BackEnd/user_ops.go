package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id" yaml:"_id" `
	UserId    string             `bson:"user_id" json:"user_id" yaml:"user_id"`
	Name      string             `json:"name" bson:"name" yaml:"name"`
	Username  string             `json:"username" bson:"username" yaml:"username"`
	Mail      string             `json:"mail" bson:"mail" yaml:"mail"`
	Password  string             `json:"password" bson:"password" yaml:"password"`
	Role      string             `json:"role" bson:"role" yaml:"role"`
	Solved    []Problems         `json:"problems_solved" bson:"problems_solved" yaml:"problems_solved"`
	Attempted []Problems         `json:"problems_attempted" bson:"problems_attempted" yaml:"problems_attempted"`
	CreatedAt int64              `json:"created_at" bson:"created_at" yaml:"created_at"`
	UpdatedAt int64              `json:"updated_at" bson:"updated_at" yaml:"updated_at"`
}

type JWTclaims struct {
	Username string
	jwt.StandardClaims
}

func show_dbs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		databases, err := client.ListDatabaseNames(ctx, bson.M{})
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err})
		}
		c.IndentedJSON(http.StatusOK, databases)
	}
}

func HashPassword(pass string) (string, error) {
	Hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(Hashed), nil
}

func GenerateToken(username string, expiry_time int64) (string, error) {
	claims := JWTclaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry_time,
		},
	}
	banana_key := "7fc76fc1fb0012d512329160df0b271075c1891816f0858b09e3e1d2c15f73b7"
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(banana_key)
	if err != nil {
		return "", err
	}
	return token, nil
}

func create_user() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var new_user User
		err := c.BindJSON(&new_user)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user_col := client.Database("BROCODE").Collection("users")
		count, err := user_col.CountDocuments(ctx, bson.M{"mail": new_user.Mail, "username": new_user.Username})
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if count != 0 {
			c.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": "user with the entered email already exists"})
			return
		} else {
			hashed_pwd, err := HashPassword(new_user.Password)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			new_user.Password = hashed_pwd
			current_time := time.Now().UnixNano()
			new_user.ID = primitive.NewObjectID()
			new_user.UserId = new_user.ID.Hex()
			new_user.Solved = []Problems{}
			new_user.Attempted = []Problems{}
			new_user.CreatedAt = current_time
			new_user.UpdatedAt = current_time
			insertionCount, err := user_col.InsertOne(ctx, new_user)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.IndentedJSON(http.StatusCreated, insertionCount)
			return
		}
	}
}

func login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		type Login struct {
			Username string `json:"username" bson:"username" yaml:"username"`
			Password string `json:"password" bson:"password" yaml:"password"`
		}
		var credentials Login
		if err := c.BindJSON(&credentials); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		hashed_pwd, err := HashPassword(credentials.Password)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		var user User
		user_col := client.Database("BROCODE").Collection("users")
		err = user_col.FindOne(ctx, bson.M{"$or": bson.A{bson.M{"username": credentials.Username}, bson.M{"mail": credentials.Username}}, "password": hashed_pwd}).Decode(&user)
		if err == mongo.ErrNoDocuments {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": "Specified credentials for login are wrong"})
			return
		}
		token_expiry := time.Now().Add(24 * time.Hour)
		token, err := GenerateToken(user.UserId, token_expiry.Unix())
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": true, "token": token})
	}
}

func get_user() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		userID := c.Param("user_id")
		var user User
		user_col := client.Database("BROCODE").Collection("users")
		err := user_col.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"status": true, "data": user})
	}
}
