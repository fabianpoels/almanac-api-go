package collections

import (
	"almanac-api/config"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetCategoryCollection(client mongo.Client) *mongo.Collection {
	return getCollection("categories", client)
}

func GetNewsItemCollection(client mongo.Client) *mongo.Collection {
	return getCollection("newsitems", client)
}

func GetUserCollection(client mongo.Client) *mongo.Collection {
	return getCollection("users", client)
}

func GetPoicollection(client mongo.Client) *mongo.Collection {
	return getCollection("pois", client)
}

func GetMunicipalityCollection(client mongo.Client) *mongo.Collection {
	return getCollection("municipalities", client)
}

func getCollection(name string, client mongo.Client) *mongo.Collection {
	return client.Database(config.GetConfig().GetString("database")).Collection(name)
}
