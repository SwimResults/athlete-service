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

var athleteCollection *mongo.Collection

func athleteService(database *mongo.Database) {
	athleteCollection = database.Collection("athlete")
}

func getAthletesByBsonDocument(d interface{}) ([]model.Athlete, error) {
	return getAthletesByBsonDocumentWithOptions(d, options.FindOptions{})
}

func getAthletesByBsonDocumentWithOptions(d interface{}, fOps options.FindOptions) ([]model.Athlete, error) {
	var athletes []model.Athlete

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := athleteCollection.Find(ctx, d, &fOps)
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

func GetAthletes(paging Paging) ([]model.Athlete, error) {
	return getAthletesByBsonDocumentWithOptions(
		bson.M{
			"$or": []interface{}{
				bson.M{"name": bson.M{"$regex": paging.Query, "$options": "i"}},
				bson.M{"firstname": bson.M{"$regex": paging.Query, "$options": "i"}},
				bson.M{"lastname": bson.M{"$regex": paging.Query, "$options": "i"}},
				bson.M{"dsv_id": bson.M{"$regex": paging.Query, "$options": "i"}},
			},
		}, paging.getPaginatedOpts())
}

func GetAthletesByMeetingId(id string, paging Paging) ([]model.Athlete, error) {
	return getAthletesByBsonDocumentWithOptions(bson.M{
		"$and": []interface{}{
			bson.M{"participation": id},
			bson.M{
				"$or": []interface{}{
					bson.M{"name": bson.M{"$regex": paging.Query, "$options": "i"}},
					bson.M{"firstname": bson.M{"$regex": paging.Query, "$options": "i"}},
					bson.M{"lastname": bson.M{"$regex": paging.Query, "$options": "i"}},
					bson.M{"dsv_id": bson.M{"$regex": paging.Query, "$options": "i"}},
				},
			},
		},
	}, paging.getPaginatedOpts())
}

func GetAthletesByTeamId(id primitive.ObjectID, paging Paging) ([]model.Athlete, error) {
	return getAthletesByBsonDocumentWithOptions(bson.M{
		"$and": []interface{}{
			bson.M{"team_id": id},
			bson.M{
				"$or": []interface{}{
					bson.M{"name": bson.M{"$regex": paging.Query, "$options": "i"}},
					bson.M{"firstname": bson.M{"$regex": paging.Query, "$options": "i"}},
					bson.M{"lastname": bson.M{"$regex": paging.Query, "$options": "i"}},
					bson.M{"dsv_id": bson.M{"$regex": paging.Query, "$options": "i"}},
				},
			},
		},
	}, paging.getPaginatedOpts())
}

func GetAthleteById(id primitive.ObjectID) (model.Athlete, error) {
	athletes, err := getAthletesByBsonDocument(bson.D{{"_id", id}})
	if err != nil {
		return model.Athlete{}, err
	}

	if len(athletes) > 0 {
		return athletes[0], nil
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

func AddParticipation(id primitive.ObjectID, meet_id string) (model.Athlete, error) {
	athlete, err := GetAthleteById(id)
	if err != nil {
		return model.Athlete{}, err
	}

	found := false
	for _, meeting := range athlete.Participation {
		if meeting == meet_id {
			found = true
		}
	}
	if !found {
		athlete.Participation = append(athlete.Participation, meet_id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = athleteCollection.ReplaceOne(ctx, bson.D{{"_id", athlete.Identifier}}, athlete)
	if err != nil {
		return model.Athlete{}, err
	}

	return GetAthleteById(athlete.Identifier)
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
