package db

import (
	"context"
	"log"
	"sci-abo-go/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"

)

var client *mongo.Client

// initialize connects to MongoDB
func InitializeDB() {

	mongoUri := config.GetEnvVar("MONGO_URI")

	// Construct the MongoDB URI string using fmt.Sprintf
	clientOptions := options.Client().ApplyURI(mongoUri)

	var err error
	client,err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// call the unique email index creation
	EnsureUserEmailUnique(client)
}

func EnsureUserEmailUnique(client *mongo.Client) {
	if client == nil {
        log.Fatal("MongoDB client is not initialized")
    }
	collection := client.Database(config.GetEnvVar("DB_NAME")).Collection("users") 

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}}, // index in ascending order
		Options: options.Index().SetUnique(true), // set the unique constraint
	}

	// create the index
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}
}

// GetCollection returns a reference to a collection in the database
func GetCollection(database string, collection string) *mongo.Collection {
	return client.Database(database).Collection(collection)
}