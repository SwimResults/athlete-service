package client

import (
	"encoding/json"
	"fmt"
	"github.com/swimresults/athlete-service/dto"
	"github.com/swimresults/athlete-service/model"
	"net/http"
)

type TeamClient struct {
	apiUrl string
}

func NewTeamClient(url string) *TeamClient {
	return &TeamClient{apiUrl: url}
}

func (c *TeamClient) importTeam(team model.Team, meeting string) (*model.Team, error) {
	request := dto.ImportTeamRequestDto{
		Meeting: meeting,
		Team:    team,
	}

	res, err := Post(c.apiUrl, "team/import", request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	newTeam := &model.Team{}
	err = json.NewDecoder(res.Body).Decode(newTeam)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("import request returned: %d", res.StatusCode)
	}
	return newTeam, nil
}
