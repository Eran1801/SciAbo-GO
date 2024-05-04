package storage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sci-abo-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

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

	// call the unique email index creation
	err = CreatingIndexesDB(client)
	if err != nil {
		log.Fatal(err)
	}

}

func CreatingIndexesDB(client *mongo.Client) error {
	if client == nil {
		return fmt.Errorf("mongoDB client is not initialized")
	}
	collection := GetUserCollection()

	index_model := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}}, // index in ascending order
		Options: options.Index().SetUnique(true),  // set the unique constraint
	}

	// create the index
	_, err := collection.Indexes().CreateOne(context.Background(), index_model)
	if err != nil {
		return fmt.Errorf("failed to create index")
	}

	log.Println("Creating indexes in mongodb was ended successfully")
	return nil
}

func GetUserCollection() *mongo.Collection {
	user_collection_name := os.Getenv("USER_COLLECTION")
	database := os.Getenv("DB_NAME")
	return client.Database(database).Collection(user_collection_name)
}

func GetUserByEmail(email string) (*models.User, error) {

	var user models.User

	collection := GetUserCollection()

	filter := bson.M{"email": email} // set the filter to retrieve data from the db

	err := collection.FindOne(context.TODO(), filter).Decode(&user) // check in the db and decode to the user variable
	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, nil // no docs was found
		}
		return nil, err
	}

	return &user, nil
}

func UpdateUser(email string, updates map[string]interface{}) error {

	collection := GetUserCollection()

	filter := bson.M{"email": email}
	update := bson.M{"$set": updates} // updates it's a map with the fields to update

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func SaveUserInDB(user *models.User, w http.ResponseWriter) error {

	// Get the User collection
	user_collection := GetUserCollection()

	// Insert the user into the MongoDB collection
	_, err := user_collection.InsertOne(context.Background(), user)
	// checks for errors
	if err != nil {
		// checks if the email is already exist in the db
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("email already exists")
		} else { // other errors
			return fmt.Errorf(err.Error())
		}
	}
	return nil
}
