package service

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sr-athlete/athlete-service/model"
	"time"
)

var teamCollection *mongo.Collection

func teamService(database *mongo.Database) {
	teamCollection = database.Collection("team")
}

func GetTeams() ([]model.Team, error) {
	var teams []model.Team

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := teamCollection.Find(ctx, bson.M{})
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

func AddTeam(team model.Team) (model.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := teamCollection.InsertOne(ctx, team)
	if err != nil {
		return model.Team{}, err
	}

	return GetTeamById(r.InsertedID.(primitive.ObjectID))
}
