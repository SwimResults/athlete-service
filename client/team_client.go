package client

import (
	"encoding/json"
	"fmt"
	"github.com/swimresults/athlete-service/dto"
	"github.com/swimresults/athlete-service/model"
	"github.com/swimresults/service-core/client"
	"net/http"
)

type TeamClient struct {
	apiUrl string
}

func NewTeamClient(url string) *TeamClient {
	return &TeamClient{apiUrl: url}
}

func (c *TeamClient) ImportTeam(team model.Team, meeting string) (*model.Team, bool, error) {
	request := dto.ImportTeamRequestDto{
		Meeting: meeting,
		Team:    team,
	}

	res, err := client.Post(c.apiUrl, "team/import", request)
	if err != nil {
		return nil, false, err
	}
	defer res.Body.Close()

	newTeam := &model.Team{}
	err = json.NewDecoder(res.Body).Decode(newTeam)
	if err != nil {
		return nil, false, err
	}

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("import request returned: %d", res.StatusCode)
	}
	return newTeam, res.StatusCode == http.StatusCreated, nil
}

func (c *TeamClient) GetTeamByName(name string) (*model.Team, bool, error) {
	fmt.Printf("request '%s'\n", c.apiUrl+"team/name?name="+name)

	params := map[string]string{
		"name": name,
	}

	res, err := client.Get(c.apiUrl, "team/name", params)
	if err != nil {
		return nil, false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("GetTeamByName received error: %d\n", res.StatusCode)
	}

	team := &model.Team{}
	err = json.NewDecoder(res.Body).Decode(team)
	if err != nil {
		return nil, false, err
	}

	return team, true, nil
}

func (c *TeamClient) GetTeamsByMeeting(meeting string) (*[]model.Team, bool, error) {
	fmt.Printf("request '%s'\n", c.apiUrl+"team/meet/"+meeting)

	res, err := client.Get(c.apiUrl, "team/meet/"+meeting, nil)
	if err != nil {
		return nil, false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("GetTeamsByMeeting received error: %d\n", res.StatusCode)
	}

	teams := &[]model.Team{}
	err = json.NewDecoder(res.Body).Decode(teams)
	if err != nil {
		return nil, false, err
	}

	return teams, true, nil
}
