package client

import (
	"fmt"
	"github.com/swimresults/athlete-service/model"
	"testing"
)

func TestImportTeam(t *testing.T) {
	team := model.Team{
		Name:    "Blubteam",
		Country: "GER",
	}
	c := NewTeamClient("http://localhost:8086/")
	newTeam, err := c.importTeam(team, "IESC12")
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Printf("id: %s, name: %s, part: %s", newTeam.Identifier.String(), newTeam.Name, newTeam.Participation)
}
