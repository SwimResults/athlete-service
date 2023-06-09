package client

import (
	"encoding/json"
	"fmt"
	"github.com/swimresults/athlete-service/dto"
	"github.com/swimresults/athlete-service/model"
	"github.com/swimresults/service-core/client"
	"net/http"
)

type AthleteClient struct {
	apiUrl string
}

func NewAthleteClient(url string) *AthleteClient {
	return &AthleteClient{apiUrl: url}
}

func (c *AthleteClient) ImportAthlete(athlete model.Athlete, meeting string) (*model.Athlete, bool, error) {
	request := dto.ImportAthleteRequestDto{
		Meeting: meeting,
		Athlete: athlete,
	}

	res, err := client.Post(c.apiUrl, "athlete/import", request)
	if err != nil {
		return nil, false, err
	}
	defer res.Body.Close()

	newAthlete := &model.Athlete{}
	err = json.NewDecoder(res.Body).Decode(newAthlete)
	if err != nil {
		return nil, false, err
	}

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("import request returned: %d", res.StatusCode)
	}
	return newAthlete, res.StatusCode == http.StatusCreated, nil
}
