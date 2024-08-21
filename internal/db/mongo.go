package db

import (
	"context"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client     *mongo.Client
	once       sync.Once
	database   *mongo.Database
)

func Connect() {
	once.Do(func() {
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			log.Fatal("Set your 'MONGODB_URI' environment variable. " +
				"See: " +
				"www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
		}
		var err error
		Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			log.Fatal(err)
		}

		err = Client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}

		dbName := os.Getenv("MONGODB_DB")
		if dbName == "" {
			log.Fatal("Set your 'MONGODB_DB' environment variable.")
		}
		database = Client.Database(dbName)
	})
}

// GetCollection returns a MongoDB collection for a given name
func GetCollection(collectionName string) *mongo.Collection {
	return database.Collection(collectionName)
}
