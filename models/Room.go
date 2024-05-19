package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {

	Content				string				`json:"content"`
	CreatedAt 			time.Time			`json:"created_at"`	
	SenderId			string				`json:"user_id_sender"`
	SenderFullName 		string				`json:"sender_full_name"`

}

// only room saved in the db
type Room struct {

	ID				primitive.ObjectID		`json:"_id" bson:"_id,omitempty"`
	User1ID 		string					`json:"user_1_id" bson:"user_1_id"`
	User2ID 		string					`json:"user_2_id" bson:"user_2_id"`
	Messages 	   []Message				`json:"messages" bson:"messages"`

}

