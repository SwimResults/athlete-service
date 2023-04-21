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

var athleteCollection *mongo.Collection

func athleteService(database *mongo.Database) {
	athleteCollection = database.Collection("athlete")
}

func GetAthletes() ([]model.Athlete, error) {
	var athletes []model.Athlete

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := athleteCollection.Find(ctx, bson.M{})
	if err != nil {
		return []model.Athlete{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var athlete model.Athlete
		cursor.Decode(&athlete)

		var team model.Team
		team, err = GetTeamById(athlete.TeamId)
		if err == nil {
			athlete.Team = team
		}

		athletes = append(athletes, athlete)
	}

	if err := cursor.Err(); err != nil {
		return []model.Athlete{}, err
	}

	return athletes, nil
}

func GetAthleteById(id primitive.ObjectID) (model.Athlete, error) {
	var athlete model.Athlete

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := athleteCollection.Find(ctx, bson.D{{"_id", id}})
	if err != nil {
		return model.Athlete{}, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		cursor.Decode(&athlete)

		var team model.Team
		team, err = GetTeamById(athlete.TeamId)
		if err == nil {
			athlete.Team = team
		}

		return athlete, nil
	}

	return model.Athlete{}, errors.New("no entry with given id found")
}

func RemoveAthleteById(id primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := athleteCollection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func AddAthlete(athlete model.Athlete) (model.Athlete, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	athlete.TeamId = athlete.Team.Identifier

	r, err := athleteCollection.InsertOne(ctx, athlete)
	if err != nil {
		return model.Athlete{}, err
	}

	return GetAthleteById(r.InsertedID.(primitive.ObjectID))
}

func UpdateAthlete(athlete model.Athlete) (model.Athlete, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	athlete.TeamId = athlete.Team.Identifier

	_, err := athleteCollection.ReplaceOne(ctx, bson.D{{"_id", athlete.Identifier}}, athlete)
	if err != nil {
		return model.Athlete{}, err
	}

	return GetAthleteById(athlete.Identifier)
}
