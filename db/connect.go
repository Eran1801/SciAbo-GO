package db

import (
	"context"
	"log"
	"sci-abo-go/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// Initialize connects to MongoDB
func InitializeDB() {
	var err error

	mongoUri := config.GetEnvVar("MONGOURI")

	// Construct the MongoDB URI string using fmt.Sprintf
	clientOptions := options.Client().ApplyURI(mongoUri)

	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
}

// GetCollection returns a reference to a collection in the database
func GetCollection(database string, collection string) *mongo.Collection {
	return Client.Database(database).Collection(collection)
}