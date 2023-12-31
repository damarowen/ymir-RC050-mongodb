// Package adapters are the glue between components and external sources.
// # This manifest was generated by ymir. DO NOT EDIT.
package adapters

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kubuskotak/ymir-test/pkg/infrastructure"
)

func TestWithUserDataMongo(t *testing.T) {
	var (
		mt      = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		adapter = &Adapter{}
	)
	defer mt.Close()

	mt.Run("connection", func(tt *mtest.T) {
		UserDataMongoOpen = func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
			return tt.Client, nil
		}
		// define env
		infrastructure.Configuration(
			infrastructure.WithPath("../.."),
			infrastructure.WithFilename("config.yaml"),
		).Initialize()

		tt.AddMockResponses(
			mtest.CreateSuccessResponse(),
		)
		adapter.Sync(
			WithUserDataMongo(&UserDataMongo{
				NetworkDB: NetworkDB{
					Host:     infrastructure.Envs.UserDataMongo.Host,
					Port:     infrastructure.Envs.UserDataMongo.Port,
					Database: infrastructure.Envs.UserDataMongo.Database,
					User:     infrastructure.Envs.UserDataMongo.User,
					Password: infrastructure.Envs.UserDataMongo.Password,
					Protocol: infrastructure.Envs.UserDataMongo.Protocol,
				},
			}),
		)
	})
}
