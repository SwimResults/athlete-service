package service

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"strings"
	"time"
)

var client *mongo.Client

func Init(c *mongo.Client) {
	database := c.Database(os.Getenv("SR_ATHLETE_MONGO_DATABASE"))
	client = c

	athleteService(database)
	teamService(database)
}

func PingDatabase() bool {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		return false
	}

	return true
}

type Paging struct {
	Limit  int
	Offset int
	Query  string
}

func (p *Paging) getPaginatedOpts() options.FindOptions {
	l := int64(p.Limit)
	skip := int64(p.Offset)
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}
	return fOpt
}

func Aliasify(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "")

	// TODO: source out in core

	return text
}

func AppendWithoutDuplicates(a []string, e string) []string {
	found := false
	for _, b := range a {
		if b == e {
			found = true
		}
	}
	if !found {
		return append(a, e)
	}
	return a
}
