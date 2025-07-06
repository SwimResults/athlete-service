package client

import (
	"encoding/json"
	"fmt"
	"github.com/swimresults/athlete-service/dto"
	"github.com/swimresults/athlete-service/model"
	"github.com/swimresults/service-core/client"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

type AthleteClient struct {
	apiUrl string
}

func NewAthleteClient(url string) *AthleteClient {
	return &AthleteClient{apiUrl: url}
}

func (c *AthleteClient) GetAthletesByMeeting(meeting string) ([]model.Athlete, error) {
	fmt.Printf("request '%s'\n", c.apiUrl+"athlete/meet/"+meeting)

	res, err := client.Get(c.apiUrl, "athlete/meet/"+meeting, nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("GetAthletesByMeeting received error: %d\n", res.StatusCode)
	}

	var athletes []model.Athlete
	err = json.NewDecoder(res.Body).Decode(&athletes)
	if err != nil {
		return nil, err
	}

	return athletes, nil
}

func (c *AthleteClient) GetAthleteByNameAndYear(name string, year int) (*model.Athlete, bool, error) {
	fmt.Printf("request '%s'\n", c.apiUrl+"athlete/name_year?name="+name+"&year="+strconv.Itoa(year))

	params := map[string]string{
		"name": name,
		"year": strconv.Itoa(year),
	}

	res, err := client.Get(c.apiUrl, "athlete/name_year", params, nil)
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

	res, err := client.Post(c.apiUrl, "athlete/import", request, nil)
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

func (c *AthleteClient) ImportCertificate(name string, athleteId primitive.ObjectID, meeting string, path string) (*model.Certificate, bool, error) {
	request := dto.ImportCertificateRequestDto{
		Name:      name,
		AthleteId: athleteId,
		Meeting:   meeting,
		Path:      path,
	}

	res, err := client.Post(c.apiUrl, "certificate/import", request, nil)
	if err != nil {
		return nil, false, err
	}
	defer res.Body.Close()

	newCertificate := &model.Certificate{}
	err = json.NewDecoder(res.Body).Decode(newCertificate)
	if err != nil {
		return nil, false, err
	}

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("import request returned: %d", res.StatusCode)
	}
	return newCertificate, res.StatusCode == http.StatusCreated, nil
}
