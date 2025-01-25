package models

import "time"

// URL represents a shortened URL with its original URL and expiration time.
type URL struct {
	// ID is the unique identifier for the shortened URL.
	ID string `json:"id" bson:"_id"`
	// URL is the original URL that has been shortened.
	URL string `json:"url" bson:"url"`
	// ExpireAt is the time when the shortened URL will expire.
	ExpireAt time.Time `json:"expireAt" bson:"expire_at"`
}
