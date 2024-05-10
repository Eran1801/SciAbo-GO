package models

import ( 
	
	"go.mongodb.org/mongo-driver/bson/primitive" 

)

type Event struct { 

	ID 					primitive.ObjectID			`json:"_id" bson:"_id,omitempty"`
	Name 				string 						`json:"name" bson:"name"`
	Abbreviation 		string 						`json:"abbreviation" bson:"abbreviation"`
	FieldOfStudy 		string 						`json:"field_of_study" bson:"field_of_study"`
	StartDate 			string						`json:"start_date" bson:"start_date"`
	EndDate 			string						`json:"end_date" bson:"end_date"`
	Country 			string						`json:"country" bson:"country"`
	City 				string 						`json:"city" bson:"city"`
}

