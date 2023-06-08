package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Team struct {
	Identifier    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`                     // automatically
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`                   // DSV-File + PDF
	Alias         []string           `json:"alias,omitempty" bson:"alias,omitempty"`                 // semiautomatically
	Country       string             `json:"country,omitempty" bson:"country,omitempty"`             // DSV-File 			!needs update!
	DsvId         int                `json:"dsv_id,omitempty" bson:"dsv_id,omitempty"`               // DSV-File 			!needs update!
	StateId       int                `json:"state_id,omitempty" bson:"state_id,omitempty"`           // DSV-File 			!needs update!
	Address       Address            `json:"address,omitempty" bson:"address,omitempty"`             // manually
	Contact       Contact            `json:"contact,omitempty" bson:"contact,omitempty"`             // manually
	Website       string             `json:"website,omitempty" bson:"website,omitempty"`             // manually
	LogoUrl       string             `json:"logo_url,omitempty" bson:"logo_url,omitempty"`           // manually
	ColorSet      ColorSet           `json:"color_set,omitempty" bson:"color_set,omitempty"`         // manually
	FirstMeeting  string             `json:"first_meeting,omitempty" bson:"first_meeting,omitempty"` // automatically
	Participation []string           `json:"participation,omitempty" bson:"participation,omitempty"` // automatically
}
