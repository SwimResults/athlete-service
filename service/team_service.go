package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/swimresults/athlete-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var teamCollection *mongo.Collection

func teamService(database *mongo.Database) {
	teamCollection = database.Collection("team")
}

func getTeamsByBsonDocument(d interface{}) ([]model.Team, error) {
	return getTeamsByBsonDocumentWithOptions(d, options.FindOptions{})
}

func getTeamsByBsonDocumentWithOptions(d interface{}, fOps options.FindOptions) ([]model.Team, error) {
	var teams []model.Team

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := teamCollection.Find(ctx, d, &fOps)
	if err != nil {
		return []model.Team{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var team model.Team
		cursor.Decode(&team)

		teams = append(teams, team)
	}

	if err := cursor.Err(); err != nil {
		return []model.Team{}, err
	}

	return teams, nil
}

func GetTeams(paging Paging) ([]model.Team, error) {
	return getTeamsByBsonDocumentWithOptions(
		bson.M{
			"$or": []interface{}{
				bson.M{"name": bson.M{"$regex": paging.Query, "$options": "i"}},
				bson.M{"alias": bson.M{"$regex": paging.Query, "$options": "i"}},
			},
		}, paging.getPaginatedOpts())
}

func GetTeamById(id primitive.ObjectID) (model.Team, error) {
	var team model.Team

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := teamCollection.Find(ctx, bson.D{{"_id", id}})
	if err != nil {
		return model.Team{}, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		cursor.Decode(&team)
		return team, nil
	}

	return model.Team{}, errors.New("no entry with given id found")
}

func GetTeamsByMeeting(id string, paging Paging) ([]model.Team, error) {
	var result = map[primitive.ObjectID]model.Team{}

	athletes, err := GetAthletesByMeetingId(id, Paging{})
	if err != nil {
		return []model.Team{}, err
	}

	for _, athlete := range athletes {
		if _, ok := result[athlete.Team.Identifier]; !ok {
			result[athlete.Team.Identifier] = athlete.Team
		}
	}

	var teams []model.Team
	for _, team := range result {
		teams = append(teams, team)
	}
	return teams, nil
}

func getTeamByName(name string) (model.Team, error) {
	teams, err := getTeamsByBsonDocument(
		bson.M{
			"$or": []interface{}{
				bson.M{"name": bson.M{"$regex": name, "$options": "i"}},
				bson.M{"alias": bson.M{"$regex": name, "$options": "i"}},
			},
		})
	if err != nil {
		return model.Team{}, err
	}
	if len(teams) < 1 {
		fmt.Printf("no team with given name '%s' found\n", name)
		return model.Team{}, errors.New("no team with given name found")
	}
	return teams[0], nil
}

func AddTeam(team model.Team) (model.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := teamCollection.InsertOne(ctx, team)
	if err != nil {
		return model.Team{}, err
	}

	return GetTeamById(r.InsertedID.(primitive.ObjectID))
}

func AddTeamParticipation(id primitive.ObjectID, meetId string) (model.Team, error) {
	fmt.Printf("add participation to team: %s (%s)\n", id.String(), meetId)
	team, err := GetTeamById(id)
	if err != nil {
		return model.Team{}, err
	}

	found := false
	for _, meeting := range team.Participation {
		if meeting == meetId {
			found = true
		}
	}
	if !found {
		team.Participation = append(team.Participation, meetId)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = teamCollection.ReplaceOne(ctx, bson.D{{"_id", team.Identifier}}, team)
	if err != nil {
		return model.Team{}, err
	}

	return GetTeamById(team.Identifier)
}

func ImportTeam(team model.Team, meetId string) (model.Team, bool, error) {
	existingTeam, err := getTeamByName(team.Name)
	if err != nil {
		if err.Error() == "no team with given name found" {
			fmt.Printf("import of team '%s', not existing so far\n", team.Name)
			newTeam, err2 := AddTeam(team)
			if err2 != nil {
				return model.Team{}, false, err2
			}
			newTeam, err2 = AddTeamParticipation(newTeam.Identifier, meetId)
			if err2 != nil {
				return model.Team{}, false, err2
			}
			return newTeam, true, nil
		}
		return model.Team{}, false, err
	}
	fmt.Printf("import of team '%s', already present\n", team.Name)
	existingTeam, err = AddTeamParticipation(existingTeam.Identifier, meetId)
	if err != nil {
		return model.Team{}, false, err
	}
	return existingTeam, false, nil
}
