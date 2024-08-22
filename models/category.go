package models

import (
	"almanac-api/config"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Category struct {
	Id     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Key    string             `json:"key" bson:"key" validate:"required,alpha"`
	Color  string             `json:"color" bson:"color" validate:"required,alpha"`
	Icon   string             `json:"icon" bson:"icon" validate:"required,alpha"`
	Active bool               `json:"active" bson:"active"`
}

func GetCategoryCollection(client mongo.Client) *mongo.Collection {
	return client.Database(config.GetConfig().GetString("database")).Collection("categories")
}
