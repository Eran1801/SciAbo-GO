package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Define a generic model interface
type Model interface{}

// initialize connects to MongoDB
func InitializeDB() {

	mongo_uri := os.Getenv("MONGO_URI")

	// construct the MongoDB URI
	client_options := options.Client().ApplyURI(mongo_uri)

	var err error
	client, err = mongo.Connect(context.TODO(), client_options)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// create unique index for User model by email
	user_collection_name := os.Getenv("USER_COLLECTION")
	index_options := options.Index().SetUnique(true)
	keys := bson.D{{Key: "email", Value: 1}}

	err = CreatingIndexesDB(user_collection_name, keys, index_options, client)
	if err != nil {
		log.Fatal(err)
	}

	// create ttl deleting after 5 minutes index in the db for PasswordResetData model
	reset_code_collection_name := os.Getenv("RESET_CODE_COLLECTION")
	index_options = options.Index().SetExpireAfterSeconds(300) // 300 seconds = 5 minutes
	keys = bson.D{{Key: "time", Value: 1}}

	err = CreatingIndexesDB(reset_code_collection_name, keys, index_options, client)
	if err != nil {
		log.Fatal(err)
	}

}

func CreatingIndexesDB(collection_name string, keys bson.D, index_options *options.IndexOptions, client *mongo.Client) error {
	if client == nil {
		return fmt.Errorf("mongoDB client is not initialized")
	}
	collection := GetCollection(collection_name)

	index_model := mongo.IndexModel{
		Keys:    keys,          // index in ascending order
		Options: index_options, // set the unique constraint
	}

	// create the index
	_, err := collection.Indexes().CreateOne(context.Background(), index_model)
	if err != nil {
		return fmt.Errorf("failed to create index")
	}

	log.Printf("Creating indexes in mongodb for %s collection was ended successfully", collection_name)
	return nil
}

func UpdateDocDB(collection_name string, id primitive.ObjectID, updates map[string]interface{}) error {

	collection := GetCollection(collection_name)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates} // updates it's a map with the fields to update

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}


func GetCollection(collection_name string) *mongo.Collection {
	database := os.Getenv("DB_NAME")
	return client.Database(database).Collection(collection_name)
}
