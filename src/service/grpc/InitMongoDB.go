package grpc

import (
	"context"
	"time"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MustInitMongoDB(config *Config, logger *shared.Logger) *mongo.Client {
	var err error
	if config.EventStoreDB != "mongodb" {
		return nil
	}
	logger.Info().Msg("bootstrapMongoDB: opening Mongodb DB connection ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Mongodb.DSN))
	if err != nil {
		logger.Panic().Msgf("bootstrapMongoDB: failed to open Mongodb DB connection: %s", err)
	}
	/***/

	logger.Info().Msg("bootstrapMongoDB: running DB migrations for customer ...")
	return client
}
