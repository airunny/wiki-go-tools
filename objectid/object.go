package objectid

import "go.mongodb.org/mongo-driver/bson/primitive"

func ObjectID() string {
	return primitive.NewObjectID().Hex()
}
