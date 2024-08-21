package db

import "go.mongodb.org/mongo-driver/mongo"

const (
	UserCollectionName = "users"
	PostCollectionName = "posts"
)

// UserCollection returns the MongoDB collection for Users
func UserCollection() *mongo.Collection {
	return GetCollection(UserCollectionName)
}

// PostCollection returns the MongoDB collection for Posts
func PostCollection() *mongo.Collection {
	return GetCollection(PostCollectionName)
}
