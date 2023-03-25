package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Athlete struct {
	Identifier primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Year       int                `json:"year,omitempty" bson:"year,omitempty"`
	Team       Team               `json:"team,omitempty" bson:"team,omitempty"` // json only
}
