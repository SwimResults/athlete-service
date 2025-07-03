package service

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/swimresults/athlete-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var certificateCollection *mongo.Collection
var certificateLogFields log.Fields

func certificateService(database *mongo.Database) {
	certificateCollection = database.Collection("certificate")
	certificateLogFields = log.Fields{"sr_service": "certificate"}
}

func getCertificatesByBsonDocument(d interface{}) ([]model.Certificate, error) {
	return getCertificatesByBsonDocumentWithOptions(d, options.FindOptions{})
}

func getCertificatesByBsonDocumentWithOptions(d interface{}, fOps options.FindOptions) ([]model.Certificate, error) {
	var certificates []model.Certificate

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fOps.SetSort(bson.D{{"ordering", 1}})

	cursor, err := certificateCollection.Find(ctx, d, &fOps)
	if err != nil {
		return []model.Certificate{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var certificate model.Certificate
		cursor.Decode(&certificate)

		certificates = append(certificates, certificate)
	}

	if err := cursor.Err(); err != nil {
		return []model.Certificate{}, err
	}

	return certificates, nil
}

func GetCertificatesAmount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Count().SetHint("_id_")
	count, err := certificateCollection.CountDocuments(ctx, bson.D{}, opts)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func GetCertificatesAmountByMeeting(meeting string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Count().SetHint("_id_")
	count, err := certificateCollection.CountDocuments(ctx, bson.D{{"meeting", meeting}}, opts)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func GetCertificates() ([]model.Certificate, error) {
	return getCertificatesByBsonDocument(bson.D{})
}

func GetCertificatesByAthleteIdAndMeeting(id primitive.ObjectID, meeting string) ([]model.Certificate, error) {
	return getCertificatesByBsonDocument(bson.D{{"athlete_id", id}, {"meeting", meeting}})
}

func GetCertificatesByAthleteId(id primitive.ObjectID) ([]model.Certificate, error) {
	return getCertificatesByBsonDocument(bson.D{{"athlete_id", id}})
}

func GetCertificateById(id primitive.ObjectID) (model.Certificate, error) {
	certificates, err := getCertificatesByBsonDocument(bson.D{{"_id", id}})
	if err != nil {
		return model.Certificate{}, err
	}

	if len(certificates) > 0 {
		return certificates[0], nil
	}

	return model.Certificate{}, errors.New("no entry with given id found")
}

func RemoveCertificateById(id primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := certificateCollection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	fields := log.Fields{"certificate_id": id}
	log.WithFields(certificateLogFields).WithFields(fields).Info("certificate deleted")
	return nil
}

func AddCertificate(certificate model.Certificate) (model.Certificate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := certificateCollection.InsertOne(ctx, certificate)
	if err != nil {
		return model.Certificate{}, err
	}

	fields := log.Fields{"certificate": certificate}
	log.WithFields(certificateLogFields).WithFields(fields).Info("certificate added")

	return GetCertificateById(r.InsertedID.(primitive.ObjectID))
}

func ImportCertificate(certificate model.Certificate) (*model.Certificate, bool, error) {

	return nil, false, nil

}

func UpdateCertificate(certificate model.Certificate) (model.Certificate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := certificateCollection.ReplaceOne(ctx, bson.D{{"_id", certificate.Identifier}}, certificate)
	if err != nil {
		return model.Certificate{}, err
	}

	fields := log.Fields{"certificate": certificate}
	log.WithFields(certificateLogFields).WithFields(fields).Info("certificate updated")

	return GetCertificateById(certificate.Identifier)
}
