package dto

import "github.com/swimresults/athlete-service/model"

type ImportAthleteRequestDto struct {
	Meeting string        `json:"meeting"`
	Athlete model.Athlete `json:"athlete"`
	Team    model.Team    `json:"team,omitempty"`
}
