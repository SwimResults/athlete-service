package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Athlete struct {
	Identifier primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Firstname  string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname   string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Year       int                `json:"year,omitempty" bson:"year,omitempty"`
	DsvId      int                `json:"dsv_id,omitempty" bson:"dsv_id,omitempty"`
	TeamId     primitive.ObjectID `json:"-" bson:"team_id,omitempty"`
	Team       Team               `json:"team,omitempty" bson:"-"`
}
