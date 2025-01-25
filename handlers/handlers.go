package handlers

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
	"url-shortener/models"
	"url-shortener/utils"
)

// Handlers struct contains the clients for MongoDB, Redis, and Snowflake node.
type Handlers struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
	node        *snowflake.Node
}

// Option is a function that configures the Handlers.
type Option func(*Handlers)

// WithMongoClient sets the MongoDB client for the Handlers.
func WithMongoClient(client *mongo.Client) Option {
	return func(h *Handlers) {
		h.mongoClient = client
	}
}

// WithRedisClient sets the Redis client for the Handlers.
func WithRedisClient(client *redis.Client) Option {
	return func(h *Handlers) {
		h.redisClient = client
	}
}

// WithNode sets the Snowflake node for the Handlers.
func WithNode(node *snowflake.Node) Option {
	return func(h *Handlers) {
		h.node = node
	}
}

// NewHandlers creates a new instance of Handlers and applies the provided options.
func NewHandlers(options ...Option) *Handlers {
	handlers := &Handlers{}

	// Apply all the functional options to configure the client.
	for _, opt := range options {
		opt(handlers)
	}

	return handlers
}

// UploadURL handles the uploading of a new URL.
// It binds the JSON payload to a URL model, checks if the long URL already exists,
// generates a unique ID based on the Snowflake node if it does not exist,
// inserts the URL into MongoDB, and returns the shortened URL.
func (h *Handlers) UploadURL(c *gin.Context) {
	var url models.URL
	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the original URL already exists in MongoDB
	collection := h.mongoClient.Database("urlshortener").Collection("urls")
	var existingURL models.URL
	if err := collection.FindOne(c, bson.M{"url": url.URL}).Decode(&existingURL); err == nil {
		// Long URL already exists, return the existing short URL
		response := map[string]string{
			"id":       existingURL.ID,
			"shortUrl": fmt.Sprintf("http://localhost:8080/%s", existingURL.ID),
		}
		c.JSON(http.StatusOK, response)
		return
	} else if !errors.Is(err, mongo.ErrNoDocuments) { // Check if the error is not "document not found"
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate unique ID and convert to base62
	id := h.node.Generate().Int64()
	url.ID = utils.Base62Encode(id)

	// Insert URL into MongoDB
	_, err := collection.InsertOne(c, url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := map[string]string{
		"id":       url.ID,
		"shortUrl": fmt.Sprintf("http://localhost:8080/%s", url.ID),
	}
	c.JSON(http.StatusOK, response)
}

// RedirectURL handles the redirection of a shortened URL to its original URL.
// It first checks the Redis cache for the shortURL URL, and if not found,
// it queries MongoDB. If the URL is found and not expired, it caches the URL in Redis
// and redirects the client to the original URL.
func (h *Handlers) RedirectURL(c *gin.Context) {
	shortURL := c.Param("id")

	// Check Redis cache first
	originalURL, err := h.redisClient.Get(c, shortURL).Result()
	if errors.Is(err, redis.Nil) {
		// If not found in cache, check MongoDB
		var url models.URL
		collection := h.mongoClient.Database("urlshortener").Collection("urls")
		err := collection.FindOne(c, bson.M{"_id": shortURL}).Decode(&url)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found", "error2": err.Error()})
			return
		}

		// Check if URL is expired
		if time.Now().After(url.ExpireAt) {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL expired"})
			return
		}

		// Cache the URL in Redis, TTL set to the time until the URL expires
		h.redisClient.Set(c, shortURL, url.URL, time.Until(url.ExpireAt))
		originalURL = url.URL
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}
