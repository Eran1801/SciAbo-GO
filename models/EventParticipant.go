package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive" 

)

type EventParticipant struct {
	
	Id 			  	primitive.ObjectID    		`json:"_id" bson:"_id,omitempty"`
	UserID        	string   					`json:"user_id" bson:"user_id"`
	ArrivalDate   	time.Time 					`json:"arrival_date" bson:"arrival_date"`
	DepartureDate 	time.Time		    	    `json:"departure_date" bson:"departure_date"`
	PosterTopic   	string    					`json:"poster_topic" bson:"poster_topic"`
}
