package models

import (
	"almanac-api/config"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewsItem struct {
	Id          primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Title       TranslatedField        `json:"title,omitempty" bson:"title,omitempty"`
	Category    string                 `json:"category" bson:"category" validate:"required,alpha"`
	Description TranslatedField        `json:"description" bson:"description"`
	Status      string                 `json:"status" bson:"status" validate:"required,alpha"`
	User        primitive.ObjectID     `json:"user,omitempty" bson:"user,omitempty"`
	GeoData     map[string]interface{} `json:"geoData" bson:"geoData"`
	Link        TranslatedField        `json:"link" bson:"link"`
	Source      string                 `json:"source" bson:"source"`
	Timestamp   time.Time              `json:"timestamp" bson:"timestamp"`
	CreatedAt   time.Time              `json:"-" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time              `json:"-" bson:"updatedAt,omitempty"`
}

func GetNewsItemCollection(client mongo.Client) *mongo.Collection {
	return client.Database(config.GetConfig().GetString("database")).Collection("newsitems")
}
