package db

import (
	"almanac-api/config"
	"almanac-api/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func GetDbClient() *mongo.Client {
	if mongoClient == nil {
		DbConnect()
	}
	return mongoClient
}

func DbConnect() {
	// mongoDB config
	username := config.GetEnv("MONGODB_USER")
	password := config.GetEnv("MONGODB_PASSW")
	host := config.GetEnv("MONGODB_HOST")
	port := config.GetEnv("MONGODB_PORT")
	mongoUrl := fmt.Sprintf("mongodb://%s:%s", host, port)
	if username != "" && password != "" {
		mongoUrl = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}
	clientOptions := options.Client().ApplyURI(mongoUrl)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Init connection
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("⛒ Connection Failed to Database")
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("⛒ Connection Failed to Database")
		log.Fatal(err)
	}

	log.Println("Connected to database")

	// (re) create indexes
	// USER INDEXES
	userEmailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	name, err := models.GetUserCollection(*client).Indexes().CreateOne(ctx, userEmailIndex)
	if err != nil {
		log.Fatal("⛒ Error creating User email index")
		log.Fatal(err)
	}
	log.Println("Created User index: " + name)

	// CATEGORY INDEXES
	categoryKeyIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "key", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	cat, err := models.GetCategoryCollection(*client).Indexes().CreateOne(ctx, categoryKeyIndex)
	if err != nil {
		log.Fatal("⛒ Error creating User email index")
		log.Fatal(err)
	}
	log.Println("Created User index: " + cat)

	mongoClient = client
}
