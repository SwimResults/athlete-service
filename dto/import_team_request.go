package dto

import "github.com/swimresults/athlete-service/model"

type ImportTeamRequestDto struct {
	Meeting string     `json:"meeting"`
	Team    model.Team `json:"team"`
}
