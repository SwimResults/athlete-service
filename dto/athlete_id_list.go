package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type AthleteIdList struct {
	Athletes []primitive.ObjectID `json:"athletes"`
}
