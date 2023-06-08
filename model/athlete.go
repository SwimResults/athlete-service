package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Athlete struct {
	Identifier    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`                     // automatically
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`                   // DSV-File + PDF
	Firstname     string             `json:"firstname,omitempty" bson:"firstname,omitempty"`         // DSV-File + (PDF) 	!needs update!
	Lastname      string             `json:"lastname,omitempty" bson:"lastname,omitempty"`           // DSV-File + (PDF) 	!needs update!
	Alias         []string           `json:"alias,omitempty" bson:"alias,omitempty"`                 // semiautomatically
	Year          int                `json:"year,omitempty" bson:"year,omitempty"`                   // DSV-File + PDF
	Gender        string             `json:"gender,omitempty" bson:"gender,omitempty"`               // DSV-File + (PDF) 	!needs update!
	DsvId         int                `json:"dsv_id,omitempty" bson:"dsv_id,omitempty"`               // DSV-File 			!needs update!
	TeamId        primitive.ObjectID `json:"-" bson:"team_id,omitempty"`                             // automatically
	Team          Team               `json:"team,omitempty" bson:"-"`                                // DSV-File + PDF
	FirstMeeting  string             `json:"first_meeting,omitempty" bson:"first_meeting,omitempty"` // automatically
	Participation []string           `json:"participation,omitempty" bson:"participation,omitempty"` // automatically
}
