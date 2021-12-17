package mongo

import (
	"context"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/rusrafkasimov/catalogs/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mgoDB struct {
	ctx    context.Context
	Client *mongo.Client
}

func InitDatabase(ctx context.Context, log promtail.Client, config *config.Configuration) (*mgoDB, error) {
	MongoHost, err := config.Get("MONGO_HOST")
	if err != nil || MongoHost == ""{
		log.Errorf("Error: can't parse mongo host. %s", err.Error())
		return nil, err
	}

	MongoUsername, err := config.Get("MONGO_USERNAME")
	if err != nil || MongoUsername == ""{
		log.Errorf("Error: can't parse mongo username. %s", err.Error())
		return nil, err
	}

	MongoPassword, err := config.Get("MONGO_PASSWORD")
	if err != nil || MongoPassword == ""{
		log.Errorf("Error: can't parse mongo pass. %s", err.Error())
		return nil, err
	}

	MongoDatabase, err := config.Get("MONGO_DATABASE")
	if err != nil || MongoDatabase == ""{
		log.Errorf("Error: can't parse mongo database. %s", err.Error())
		return nil, err
	}

	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    MongoDatabase,
		Username:      MongoUsername,
		Password:      MongoPassword,
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + MongoHost).SetAuth(credential)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Errorf("Error: can't connect mongo. %s", err.Error())
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Errorf("Error: can't ping mongo client. %s", err.Error())
		return nil, err
	}

	return &mgoDB{
		ctx:    ctx,
		Client: client,
	}, nil
}
