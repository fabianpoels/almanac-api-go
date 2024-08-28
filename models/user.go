package models

import (
	"almanac-api/config"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty" validate:"required,alpha"`
	Password  string             `json:"-" bson:"password,omitempty" validate:"required,alpha"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty" validate:"required,alpha"`
	Active    bool               `json:"active" bson:"active"`
	Role      string             `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"-" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"-" bson:"updatedAt,omitempty"`
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin" || u.Role == "superadmin"
}

func GetUserCollection(client mongo.Client) *mongo.Collection {
	return client.Database(config.GetConfig().GetString("database")).Collection("users")
}
