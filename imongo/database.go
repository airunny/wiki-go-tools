package imongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	*mongo.Database
	debug bool
}

func (d *Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	return &Collection{
		Collection: d.Database.Collection(name, opts...),
		debug:      d.debug,
	}
}
