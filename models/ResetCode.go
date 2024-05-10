package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResetCode struct {
	ID   	primitive.ObjectID 			`json:"id" bson:"_id,omitempty"`
	Code 	string           			`json:"code" bson:"code"`
	Time 	time.Time         			`json:"time" bson:"time"`
}
