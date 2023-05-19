package service

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sr-athlete/athlete-service/model"
	"time"
)

var teamCollection *mongo.Collection

func teamService(database *mongo.Database) {
	teamCollection = database.Collection("team")
}

func getTeamsByBsonDocument(d primitive.D) ([]model.Team, error) {
	return getTeamsByBsonDocumentWithOptions(d, options.FindOptions{})
}

func getTeamsByBsonDocumentWithOptions(d primitive.D, fOps options.FindOptions) ([]model.Team, error) {
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
	return getTeamsByBsonDocumentWithOptions(bson.D{}, paging.getPaginatedOpts())
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

func AddTeam(team model.Team) (model.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := teamCollection.InsertOne(ctx, team)
	if err != nil {
		return model.Team{}, err
	}

	return GetTeamById(r.InsertedID.(primitive.ObjectID))
}
