package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"sync"

	"github.com/mahdi-eth/social-media-graphql/internal/db"
	"github.com/mahdi-eth/social-media-graphql/api/graphql/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var postSubscriptions = make(map[string]chan *model.Post)
var mu sync.Mutex

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	userCollection := db.UserCollection()

	newUser := bson.M{
		"name":      input.Name,
		"following": []primitive.ObjectID{},
		"followers": []primitive.ObjectID{},
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	objectID := result.InsertedID.(primitive.ObjectID).Hex()

	return &model.User{
		ID:        objectID,
		Name:      *input.Name,
		Following: []*model.User{},
		Followers: []*model.User{},
	}, nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	postCollection := db.PostCollection()
	userCollection := db.UserCollection()

	authorObjectID, err := primitive.ObjectIDFromHex(input.AuthorID)
	if err != nil {
		return nil, errors.New("invalid author ID format")
	}

	var author struct {
		ID   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": authorObjectID}).Decode(&author)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	newPost := bson.M{
		"author":  authorObjectID,
		"content": input.Content,
	}

	result, err := postCollection.InsertOne(ctx, newPost)
	if err != nil {
		return nil, err
	}

	postID := result.InsertedID.(primitive.ObjectID).Hex()

	createdPost := &model.Post{
		ID: postID,
		Author: &model.User{
			ID:   author.ID.Hex(),
			Name: author.Name,
		}, Content: input.Content,
	}

	notifyPostAdded(createdPost)

	return createdPost, nil
}

// FollowUser is the resolver for the followUser field.
func (r *mutationResolver) FollowUser(ctx context.Context, followerID string, followeeID string) (*model.User, error) {
	userCollection := db.UserCollection()

	followerObjectID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		return nil, errors.New("invalid follower ID format")
	}

	followeeObjectID, err := primitive.ObjectIDFromHex(followeeID)
	if err != nil {
		return nil, errors.New("invalid followee ID format")
	}

	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": followerObjectID},
		bson.M{"$addToSet": bson.M{"following": followeeObjectID}},
	)
	if err != nil {
		return nil, err
	}

	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": followeeObjectID},
		bson.M{"$addToSet": bson.M{"followers": followerObjectID}},
	)
	if err != nil {
		return nil, err
	}

	var mongoUser struct {
		ID        primitive.ObjectID   `bson:"_id"`
		Name      string               `bson:"name"`
		Following []primitive.ObjectID `bson:"following,omitempty"`
		Followers []primitive.ObjectID `bson:"followers,omitempty"`
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": followeeObjectID}).Decode(&mongoUser)
	if err != nil {
		return nil, err
	}

	updatedFollower := &model.User{
		ID:   mongoUser.ID.Hex(),
		Name: mongoUser.Name,
	}

	return updatedFollower, nil
}

// UnfollowUser is the resolver for the unfollowUser field.
func (r *mutationResolver) UnfollowUser(ctx context.Context, followerID string, followeeID string) (*model.User, error) {
	userCollection := db.UserCollection()

	followerObjectID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		return nil, errors.New("invalid follower ID format")
	}

	followeeObjectID, err := primitive.ObjectIDFromHex(followeeID)
	if err != nil {
		return nil, errors.New("invalid followee ID format")
	}

	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": followerObjectID},
		bson.M{"$pull": bson.M{"following": followeeObjectID}},
	)
	if err != nil {
		return nil, err
	}

	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": followeeObjectID},
		bson.M{"$pull": bson.M{"followers": followerObjectID}},
	)
	if err != nil {
		return nil, err
	}
	var mongoUser struct {
		ID        primitive.ObjectID `bson:"_id"`
		Name      string             `bson:"name"`
		Following []*model.User      `bson:"following,omitempty"`
		Followers []*model.User      `bson:"followers,omitempty"`
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": followerObjectID}).Decode(&mongoUser)
	if err != nil {
		return nil, err
	}

	updatedFollower := &model.User{
		ID:        mongoUser.ID.Hex(),
		Name:      mongoUser.Name,
		Following: mongoUser.Following,
		Followers: mongoUser.Followers,
	}

	return updatedFollower, nil
}

// PostsByFollowing is the resolver for the postsByFollowing field.
func (r *queryResolver) PostsByFollowing(ctx context.Context, userID string) ([]*model.Post, error) {
	userCollection := db.UserCollection()
	postCollection := db.PostCollection()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var user struct {
		Following []primitive.ObjectID `bson:"following,omitempty"`
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	if len(user.Following) == 0 {
		return []*model.Post{}, nil
	}

	cursor, err := postCollection.Find(ctx, bson.M{"author": bson.M{"$in": user.Following}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*model.Post
	for cursor.Next(ctx) {
		var post model.Post
		var mongoPost struct {
			ID      primitive.ObjectID `bson:"_id"`
			Author  primitive.ObjectID `bson:"author"`
			Content string             `bson:"content"`
		}

		if err := cursor.Decode(&mongoPost); err != nil {
			return nil, err
		}

		post.ID = mongoPost.ID.Hex()
		post.Content = mongoPost.Content
		post.Author = &model.User{
			ID: mongoPost.Author.Hex(),
		}

		posts = append(posts, &post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// PostAddedByFollowing is the resolver for the postAddedByFollowing field.
func (r *subscriptionResolver) PostAddedByFollowing(ctx context.Context, userID string) (<-chan *model.Post, error) {
	userCollection := db.UserCollection()
	postCollection := db.PostCollection()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var user struct {
		Following []primitive.ObjectID `bson:"following,omitempty"`
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	postChan := make(chan *model.Post)

	mu.Lock()
	postSubscriptions[userID] = postChan
	mu.Unlock()

	go func() {
		defer close(postChan)

		cursor, err := postCollection.Find(ctx, bson.M{"author": bson.M{"$in": user.Following}})
		if err != nil {
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var mongoPost struct {
				ID      primitive.ObjectID `bson:"_id"`
				Author  primitive.ObjectID `bson:"author"`
				Content string             `bson:"content"`
			}

			if err := cursor.Decode(&mongoPost); err != nil {
				continue
			}

			var mongoUser struct {
				ID        primitive.ObjectID `bson:"_id"`
				Name      string             `bson:"name"`
			}
		
			err = userCollection.FindOne(ctx, bson.M{"_id": mongoPost.Author}).Decode(&mongoUser)
			if err != nil {
				continue
			}

			post := &model.Post{
				ID: mongoPost.ID.Hex(),
				Content: mongoPost.Content,
				Author: &model.User{
					ID: mongoPost.Author.Hex(),
					Name: mongoUser.Name,
				},
			}

			postChan <- post
		}

		for {
			select {
			case <-ctx.Done():
				mu.Lock()
				delete(postSubscriptions, userID)
				mu.Unlock()
				return
			}
		}
	}()

	return postChan, nil
}


// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }