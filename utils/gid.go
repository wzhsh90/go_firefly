package utils

import "github.com/globalsign/mgo/bson"

func NewId() string {
	return bson.NewObjectId().Hex()
}
