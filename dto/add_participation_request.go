package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddParticipationRequestDto struct {
	AthleteId primitive.ObjectID `json:"athlete,omitempty"`
	MeetingId string             `json:"meeting,omitempty"`
}
