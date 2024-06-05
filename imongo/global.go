package imongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	defaultClient *Client
)

func GetDatabase(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	if defaultClient == nil {
		panic("nil mongo client")
	}

	return defaultClient.Database(name, opts...)
}

func GetIDatabase(name string, opts ...*options.DatabaseOptions) *Database {
	if defaultClient == nil {
		panic("nil mongo client")
	}

	return defaultClient.IDatabase(name, opts...)
}

func IsMongoDuplicate(err error) bool {
	var e mongo.WriteException
	ok := errors.As(err, &e)
	if !ok {
		return false
	}

	if len(e.WriteErrors) > 0 && e.WriteErrors[0].Code == 11000 {
		return true
	}
	return false
}

func StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	if defaultClient == nil {
		panic("nil mongo client")
	}

	return defaultClient.StartSession(opts...)
}

func Close() error {
	if defaultClient == nil {
		return nil
	}
	return defaultClient.Disconnect(context.Background())
}
