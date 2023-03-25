package service

import (
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func Init(c *mongo.Client) {
	database := c.Database(os.Getenv("SR_ATHLETE_MONGO_DATABASE"))

	athleteService(database)
	teamService(database)
}
