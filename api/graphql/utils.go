package graph

import (
	"context"
	"errors"
	"fmt"

	"github.com/mahdi-eth/social-media-graphql/api/graphql/model"
	"github.com/mahdi-eth/social-media-graphql/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// This function will be called when a new post is added
func notifyPostAdded(post *model.Post) {
	mu.Lock()
	defer mu.Unlock()

	for userID, postChan := range postSubscriptions {
		isFollowing, err := isUserFollowingAuthor(userID, post.Author.ID)
		if err != nil {
			fmt.Printf("Error checking if user %s is following author %s: %v\n", userID, post.Author.ID, err)
			continue
		}

		if isFollowing {
			select {
				case postChan <- post:
				default:
					fmt.Printf("Post channel for user %s is full, skipping notification\n", userID)
			}
		}
	}
}


func isUserFollowingAuthor(userID, authorID string) (bool, error) {
    userCollection := db.UserCollection()

    userObjectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return false, errors.New("invalid user ID format")
    }

    authorObjectID, err := primitive.ObjectIDFromHex(authorID)
    if err != nil {
        return false, errors.New("invalid author ID format")
    }

    var user struct {
        Following []primitive.ObjectID `bson:"following,omitempty"`
    }

    err = userCollection.FindOne(context.TODO(), bson.M{"_id": userObjectID}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return false, nil
        }
        return false, err
    }

    for _, followedID := range user.Following {
        if followedID == authorObjectID {
            return true, nil
        }
    }

    return false, nil
}
