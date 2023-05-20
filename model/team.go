package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Team struct {
	Identifier primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Alias      []string           `json:"alias,omitempty" bson:"alias,omitempty"`
	Address    Address            `json:"address,omitempty" bson:"address,omitempty"`
	FirstEvent string             `json:"first_event,omitempty" bson:"first_event,omitempty"`
	Country    string             `json:"country,omitempty" bson:"country,omitempty"`
	Contact    Contact            `json:"contact,omitempty" bson:"contact,omitempty"`
	Website    string             `json:"website,omitempty" bson:"website,omitempty"`
	LogoUrl    string             `json:"logo_url,omitempty" bson:"logo_url,omitempty"`
	ColorSet   ColorSet           `json:"color_set,omitempty" bson:"color_set,omitempty"`
}
