package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct { 

	ID 					primitive.ObjectID			`json:"_id" bson:"_id,omitempty"` 
	Name 				string 						`json:"name" bson:"name"`
	Abbreviation 		string 						`json:"abbreviation" bson:"abbreviation"`
	Field				string 						`json:"field" bson:"field"`
	StartDate 			time.Time						`json:"start_date" bson:"start_date"`
	EndDate 			time.Time						`json:"end_date" bson:"end_date"`
	Country 			string						`json:"country" bson:"country"`
	City 				string 						`json:"city" bson:"city"`
	Participants	  []string						`json:"participants" bson:"participants"`

}

