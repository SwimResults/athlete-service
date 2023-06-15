package client

import (
	"encoding/json"
	"fmt"
	"github.com/swimresults/athlete-service/dto"
	"github.com/swimresults/athlete-service/model"
	"github.com/swimresults/service-core/client"
	"net/http"
	"strconv"
)

type AthleteClient struct {
	apiUrl string
}

func NewAthleteClient(url string) *AthleteClient {
	return &AthleteClient{apiUrl: url}
}

func (c *AthleteClient) GetAthleteByNameAndYear(name string, year int) (*model.Athlete, bool, error) {
	fmt.Printf("request '%s'\n", c.apiUrl+"athlete/name_year?name="+name+"&year="+strconv.Itoa(year))

	params := map[string]string{
		"name": name,
		"year": strconv.Itoa(year),
	}

	res, err := client.Get(c.apiUrl, "athlete/name_year", params)
	if err != nil {
		return nil, false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("GetAthleteByNameAndYear received error: %d\n", res.StatusCode)
	}

	athlete := &model.Athlete{}
	err = json.NewDecoder(res.Body).Decode(athlete)
	if err != nil {
		return nil, false, err
	}

	return athlete, true, nil
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
