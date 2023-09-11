package mongodb

import (
	"context"
	"fmt"
	"vitalsign-publisher/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	conf = config.GetConfig()
)

func GetMongoClient() (context.Context, context.CancelFunc, *mongo.Client, error) {
	// Database connection
	mongoURL := fmt.Sprintf("mongodb://%s:%d", conf.MongoDB.Host, conf.MongoDB.Port)
	ctx, cancel := context.WithCancel(context.Background())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return ctx, cancel, nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return ctx, cancel, nil, err
	}
	return ctx, cancel, client, nil
}
