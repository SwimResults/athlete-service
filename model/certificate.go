package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Certificate struct {
	Identifier primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	AthleteId  primitive.ObjectID `json:"athlete_id,omitempty" bson:"athlete_id,omitempty"`
	Meeting    string             `json:"meeting,omitempty" bson:"meeting,omitempty"`
	Path       string             `json:"path,omitempty" bson:"path,omitempty"`
	Url        string             `json:"url,omitempty" bson:"url,omitempty"`
	Hidden     bool               `json:"hidden,omitempty" bson:"hidden,omitempty"`
	Downloads  int                `json:"downloads,omitempty" bson:"downloads,omitempty"`
	Ordering   int                `json:"ordering,omitempty" bson:"ordering,omitempty"`
	AddedAt    time.Time          `json:"added_at,omitempty" bson:"added_at,omitempty"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
