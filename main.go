package main

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"url-shortener/handlers"
)

var (
	mongoClient *mongo.Client
	redisClient *redis.Client
	node        *snowflake.Node
)

func main() {
	var err error

	// Using Snowflake to generate unique IDs which will ensure the id would be unique across all instances of the service.
	//
	// @NOTE: Initialize Snowflake node with a node number of 1
	// @NOTE: if we have multiple instances of the service running, we need to assign a unique node number to each instance to avoid id conflicts.
	node, err = snowflake.NewNode(1)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize MongoDB client
	ctx := context.Background()
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(ctx)

	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	defer redisClient.Close()

	r := gin.Default()
	h := handlers.NewHandlers(
		handlers.WithMongoClient(mongoClient),
		handlers.WithRedisClient(redisClient),
		handlers.WithNode(node),
	)
	r.POST("/api/v1/urls", h.UploadURL)
	r.GET("/:id", h.RedirectURL)

	r.Run(":8080")
}
