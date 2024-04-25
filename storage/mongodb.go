package storage

import (
	"context"
	"log"
	"sci-abo-go/config"
	"sci-abo-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// initialize connects to MongoDB
func InitializeDB() {

	mongoUri := config.GetEnvVar("MONGO_URI")

	// Construct the MongoDB URI string using fmt.Sprintf
	clientOptions := options.Client().ApplyURI(mongoUri)

	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// call the unique email index creation
	CreatingIndexesByUserEmail(client)
}

func CreatingIndexesByUserEmail(client *mongo.Client) {
	if client == nil {
		log.Fatal("MongoDB client is not initialized")
	}
	collection := GetCollection(config.GetEnvVar("USER_COLLECTION"))

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}}, // index in ascending order
		Options: options.Index().SetUnique(true),  // set the unique constraint
	}

	// create the index
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}
}

// GetCollection returns a reference to a collection in the database
func GetCollection(collection string) *mongo.Collection {
	database := config.GetEnvVar("DB_NAME")
	return client.Database(database).Collection(collection)
}

func GetUserById(id string) (*models.User, error) {

	var user models.User

	collection := GetCollection(config.GetEnvVar("USER_COLLECTION"))
	filter := bson.M{"id": id}

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, nil // no docs was found
		}
		return nil, err
	}

	return &user, nil
}

func UpdateUser(id string, updates map[string]interface{}) error {
	collection := GetCollection(config.GetEnvVar("USER_COLLECTION"))
	filter := bson.M{"id": id}
	update := bson.M{"$set": updates}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
