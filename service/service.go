package service

import (
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func Init(c *mongo.Client) {
	database := c.Database(os.Getenv("SR_EXAMPLE_MONGO_DATABASE"))

	exampleService(database)
}
