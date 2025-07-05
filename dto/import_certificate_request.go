package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImportCertificateRequestDto struct {
	Name      string             `json:"name"`
	Meeting   string             `json:"meeting"`
	AthleteId primitive.ObjectID `json:"athlete_id"`
	Path      string             `json:"path"`
}
