package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/swimresults/athlete-service/model"
	"github.com/swimresults/service-core/misc"
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

	fOps.SetSort(bson.D{{"name", 1}})

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
				bson.M{"alias": bson.M{"$regex": misc.Aliasify(paging.Query), "$options": "i"}},
			},
		}, paging.getPaginatedOpts())
}

func GetTeamsByMeeting(id string, paging Paging) ([]model.Team, error) {
	return getTeamsByBsonDocumentWithOptions(bson.M{
		"$and": []interface{}{
			bson.M{"participation": id},
			bson.M{
				"$or": []interface{}{
					bson.M{"name": bson.M{"$regex": paging.Query, "$options": "i"}},
					bson.M{"alias": bson.M{"$regex": paging.Query, "$options": "i"}},
					bson.M{"alias": bson.M{"$regex": misc.Aliasify(paging.Query), "$options": "i"}},
				},
			},
		},
	}, paging.getPaginatedOpts())
}

func GetTeamById(id primitive.ObjectID) (model.Team, error) {
	teams, err := getTeamsByBsonDocument(bson.D{{"_id", id}})
	if err != nil {
		return model.Team{}, err
	}
	if len(teams) < 1 {
		fmt.Printf("no team with given id '%d' found\n", id)
		return model.Team{}, errors.New("no entry with given id found")
	}
	return teams[0], nil
}

func GetTeamByDsvId(dsvId int) (model.Team, error) {
	teams, err := getTeamsByBsonDocument(bson.D{{"dsv_id", dsvId}})
	if err != nil {
		return model.Team{}, err
	}
	if len(teams) < 1 {
		fmt.Printf("no team with given dsv_id '%d' found\n", dsvId)
		return model.Team{}, errors.New("no entry with given id found")
	}
	return teams[0], nil
}

func GetTeamByName(name string) (model.Team, error) {
	teams, err := getTeamsByBsonDocument(
		bson.M{
			"$or": []interface{}{
				bson.M{"name": bson.M{"$regex": name, "$options": "i"}},
				bson.M{"alias": bson.M{"$regex": name, "$options": "i"}},
				bson.M{"alias": bson.M{"$regex": misc.Aliasify(name), "$options": "i"}},
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

func GetTeamByAlias(alias string) (model.Team, error) {
	teams, err := getTeamsByBsonDocument(
		bson.M{
			"$or": []interface{}{
				bson.M{"name": bson.M{"$regex": alias, "$options": "i"}},
				bson.M{"alias": bson.M{"$regex": alias, "$options": "i"}},
			},
		})
	if err != nil {
		return model.Team{}, err
	}
	if len(teams) < 1 {
		fmt.Printf("no team with given alias '%s' found\n", alias)
		return model.Team{}, errors.New("no team with given alias found")
	}
	return teams[0], nil
}

func AddTeam(team model.Team) (model.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	team.Alias = append(team.Alias, misc.Aliasify(team.Name))

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

	team.Participation = misc.AppendWithoutDuplicates(team.Participation, meetId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = teamCollection.ReplaceOne(ctx, bson.D{{"_id", team.Identifier}}, team)
	if err != nil {
		return model.Team{}, err
	}

	return GetTeamById(team.Identifier)
}

func ImportTeam(team model.Team, meetId string) (model.Team, bool, error) {
	existingTeam, err := GetTeamByName(team.Name)

	if err != nil {
		if err.Error() == "no team with given name found" {

			fmt.Printf("import of team '%s', not existing so far\n", team.Name)
			team.FirstMeeting = meetId
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

	changed := false
	if existingTeam.DsvId == 0 && team.DsvId != 0 {
		existingTeam.DsvId = team.DsvId
		changed = true
	}
	if existingTeam.StateId == 0 && team.StateId != 0 {
		existingTeam.StateId = team.StateId
		changed = true
	}
	if existingTeam.Country == "" && team.Country != "" {
		existingTeam.Country = team.Country
		changed = true
	}

	fmt.Printf("import of team '%s', already present\n", team.Name)

	if changed {
		fmt.Printf("updating some values...\n")
		existingTeam, err = UpdateTeam(existingTeam)
		if err != nil {
			return model.Team{}, false, err
		}
	}

	existingTeam, err = AddTeamParticipation(existingTeam.Identifier, meetId)
	if err != nil {
		return model.Team{}, false, err
	}

	return existingTeam, false, nil
}

func UpdateTeam(team model.Team) (model.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	team.Alias = misc.AppendWithoutDuplicates(team.Alias, misc.Aliasify(team.Name))

	_, err := teamCollection.ReplaceOne(ctx, bson.D{{"_id", team.Identifier}}, team)
	if err != nil {
		return model.Team{}, err
	}

	return GetTeamById(team.Identifier)
}
