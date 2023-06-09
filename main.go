package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/swimresults/athlete-service/controller"
	"github.com/swimresults/athlete-service/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"math/rand"
	"os"
	"time"
)

var client *mongo.Client

func main() {
	ctx := connectDB()

	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime: "@timestamp",
			log.FieldKeyMsg:  "message",
		},
	})
	log.SetLevel(log.TraceLevel)

	min := 1000000
	max := 9999999
	rnd := rand.Intn(max-min) + min
	filename := fmt.Sprintf("logs/out-%d.log", rnd)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	}
	defer file.Close()

	log.Info("athlete service started")

	service.Init(client)
	controller.Run()

	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func connectDB() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	var uri = "mongodb://"
	if os.Getenv("SR_ATHLETE_MONGO_USERNAME") != "" {
		uri += os.Getenv("SR_ATHLETE_MONGO_USERNAME") + ":" + os.Getenv("SR_ATHLETE_MONGO_PASSWORD") + "@"
	}
	uri += os.Getenv("SR_ATHLETE_MONGO_HOST") + ":" + os.Getenv("SR_ATHLETE_MONGO_PORT")
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		fmt.Println("failed when trying to connect to '" + os.Getenv("SR_ATHLETE_MONGO_HOST") + ":" + os.Getenv("SR_ATHLETE_MONGO_PORT") + "' as '" + os.Getenv("SR_ATHLETE_MONGO_USERNAME") + "'")
		fmt.Println(fmt.Errorf("unable to connect to mongo database"))
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("failed when trying to connect to '" + os.Getenv("SR_ATHLETE_MONGO_HOST") + ":" + os.Getenv("SR_ATHLETE_MONGO_PORT") + "' as '" + os.Getenv("SR_ATHLETE_MONGO_USERNAME") + "'")
		fmt.Println(fmt.Errorf("unable to reach mongo database"))
	}

	return ctx
}
