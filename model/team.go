package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Team struct {
	Identifier primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Alias      []string           `json:"alias,omitempty" bson:"alias,omitempty"`
	Address    Address            `json:"address,omitempty" bson:"address,omitempty"`
}
