package service

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/swimresults/athlete-service/model"
	"github.com/swimresults/service-core/misc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var athleteCollection *mongo.Collection
var athleteLogFields log.Fields

func athleteService(database *mongo.Database) {
	athleteCollection = database.Collection("athlete")
	athleteLogFields = log.Fields{"sr_service": "athlete"}
}

func getAthletesByBsonDocument(d interface{}) ([]model.Athlete, error) {
	return getAthletesByBsonDocumentWithOptions(d, options.FindOptions{})
}

func getAthletesByBsonDocumentWithOptions(d interface{}, fOps options.FindOptions) ([]model.Athlete, error) {
	var athletes []model.Athlete

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fOps.SetSort(bson.D{{"name", 1}})

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
				bson.M{"alias": bson.M{"$regex": paging.Query, "$options": "i"}},
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
					bson.M{"alias": bson.M{"$regex": paging.Query, "$options": "i"}},
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
					bson.M{"alias": bson.M{"$regex": paging.Query, "$options": "i"}},
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

func GetAthleteByDsvId(dsvId int) (model.Athlete, error) {
	athletes, err := getAthletesByBsonDocument(bson.D{{"dsv_id", dsvId}})
	if err != nil {
		return model.Athlete{}, err
	}

	if len(athletes) > 0 {
		return athletes[0], nil
	}

	return model.Athlete{}, errors.New("no entry with given dsv_id found")
}

func GetAthleteByNameAndYear(name string, year int) (model.Athlete, error) {
	if hasComma, first, last := extractNames(name); hasComma {
		name = first + " " + last
	}

	athletes, err := getAthletesByBsonDocument(bson.M{
		"$and": []interface{}{
			bson.M{"year": year},
			bson.M{
				"$or": []interface{}{
					bson.M{"name": bson.M{"$regex": name, "$options": "i"}},
					bson.M{"alias": bson.M{"$regex": misc.Aliasify(name), "$options": "i"}},
				},
			},
		},
	})
	if err != nil {
		return model.Athlete{}, err
	}

	if len(athletes) > 0 {
		return athletes[0], nil
	}

	fmt.Printf("no athlete with given name '%s' and year %d found\n", name, year)
	return model.Athlete{}, errors.New("no entry found")
}

func RemoveAthleteById(id primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := athleteCollection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	fields := log.Fields{"athlete_id": id}
	log.WithFields(athleteLogFields).WithFields(fields).Info("athlete deleted")
	return nil
}

func AddAthlete(athlete model.Athlete) (model.Athlete, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	athlete.TeamId = athlete.Team.Identifier

	if hasComma, first, last := extractNames(athlete.Name); hasComma {
		athlete.Name = first + " " + last
		athlete.Firstname = first
		athlete.Lastname = last
	}

	athlete.Alias = misc.AppendWithoutDuplicates(athlete.Alias, misc.Aliasify(athlete.Name))

	r, err := athleteCollection.InsertOne(ctx, athlete)
	if err != nil {
		return model.Athlete{}, err
	}

	fields := log.Fields{"athlete": athlete}
	log.WithFields(athleteLogFields).WithFields(fields).Info("athlete added")

	return GetAthleteById(r.InsertedID.(primitive.ObjectID))
}

func AddParticipation(id primitive.ObjectID, meetId string) (model.Athlete, error) {
	fmt.Printf("add participation to athlete: %s (%s)\n", id.String(), meetId)
	athlete, err := GetAthleteById(id)
	if err != nil {
		return model.Athlete{}, err
	}

	athlete.Participation = misc.AppendWithoutDuplicates(athlete.Participation, meetId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = athleteCollection.ReplaceOne(ctx, bson.D{{"_id", athlete.Identifier}}, athlete)
	if err != nil {
		return model.Athlete{}, err
	}

	fields := log.Fields{"athlete": athlete, "meet_id": meetId}
	log.WithFields(athleteLogFields).WithFields(fields).Info("participation added to athlete")

	return GetAthleteById(athlete.Identifier)
}

func ImportAthlete(athlete model.Athlete, meetId string) (*model.Athlete, bool, error) {

	if athlete.Team.Name == "" && athlete.Team.DsvId == 0 {
		return nil, false, fmt.Errorf("no team set in import")
	}

	var existing model.Athlete
	var err error
	found := false
	if athlete.DsvId != 0 {
		existing, err = GetAthleteByDsvId(athlete.DsvId)

		if err != nil {
			if err.Error() != "no entry with given dsv_id found" {
				return nil, false, err
			}
		} else {
			found = true
		}
	}

	if !found {
		existing, err = GetAthleteByNameAndYear(athlete.Name, athlete.Year)

		if err != nil {
			if err.Error() != "no entry found" {
				return nil, false, err
			}
		} else {
			found = true
		}
	}

	if found {
		fmt.Printf("import of athlete '%s', already present\n", athlete.Name)

		changed := false
		if existing.Firstname == "" || existing.Lastname == "" {
			if hasNames, first, last := extractNames(athlete.Name); hasNames {
				existing.Firstname = first
				existing.Lastname = last
				changed = true
			}
		}
		if existing.DsvId == 0 && athlete.DsvId != 0 {
			existing.DsvId = athlete.DsvId
			changed = true
		}
		if existing.Gender == "" && athlete.Gender != "" {
			existing.Gender = athlete.Gender
			changed = true
		}

		if changed {
			fmt.Printf("updating some values...\n")
			existing, err = UpdateAthlete(existing)
			if err != nil {
				return nil, false, err
			}
		}
	} else {
		fmt.Printf("import of athlete '%s', not existing so far\n", athlete.Name)

		athlete.FirstMeeting = meetId

		var team model.Team
		if athlete.Team.DsvId != 0 {
			team, err = GetTeamByDsvId(athlete.Team.DsvId)
		} else {
			team, err = getTeamByName(athlete.Team.Name)
		}
		if err != nil {
			return nil, true, err
		}

		athlete.Team.Identifier = team.Identifier

		existing, err = AddAthlete(athlete)
		if err != nil {
			return nil, true, err
		}
	}

	existing, err = AddParticipation(existing.Identifier, meetId)
	if err != nil {
		return nil, !found, err
	}

	fields := log.Fields{"athlete": existing, "created": !found}
	log.WithFields(athleteLogFields).WithFields(fields).Info("athlete imported")

	return &existing, !found, nil

	// if dsv_id, search by dsv_id (dsv_id '==')
	// -> not found
	// 		search by name, team and year (name aliasified in aliases; team (getTeamByName), year '==')
	//
	// 		-> not found:	create
	//						add participation
	// 			=> return true
	//
	// -> found:	update dsv_id
	// 				update firstname (+aliases)
	//				update lastname (+aliases)
	// 				update gender
	//				add participation
	// 			=> return false

}

func UpdateAthlete(athlete model.Athlete) (model.Athlete, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	athlete.TeamId = athlete.Team.Identifier

	if hasComma, first, last := extractNames(athlete.Name); hasComma {
		athlete.Name = first + " " + last
		athlete.Firstname = first
		athlete.Lastname = last
	}

	athlete.Alias = misc.AppendWithoutDuplicates(athlete.Alias, misc.Aliasify(athlete.Name))

	_, err := athleteCollection.ReplaceOne(ctx, bson.D{{"_id", athlete.Identifier}}, athlete)
	if err != nil {
		return model.Athlete{}, err
	}

	fields := log.Fields{"athlete": athlete}
	log.WithFields(athleteLogFields).WithFields(fields).Info("athlete updated")

	return GetAthleteById(athlete.Identifier)
}
