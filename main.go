package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"almanac-api/collections"
	"almanac-api/config"
	"almanac-api/db"
	"almanac-api/server"
	"almanac-api/utils"

	"gitlab.com/almanac-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	env := config.GetEnv("ENVIRONMENT")

	environment := flag.String("e", env, "")
	migration := flag.Bool("migrate", false, "Run db migration task")

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	db.DbConnect()
	db.CacheConnect()

	if *migration {
		log.Println("MIGRATING")
		riskLevelsMigration()
		return
	}

	// start server
	server.Init()
}

func importLocations() {
	file, err := os.Open("locations_coordinates.csv")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','          // Set delimiter (default is comma)
	reader.FieldsPerRecord = -1 // -1 means no validation on number of fields

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	for i, record := range records {
		if i == 0 {
			fmt.Println("Headers:", record)
			continue
		}

		coords := strings.Split(record[2], ", ")
		lat, _ := strconv.ParseFloat(coords[1], 64)
		long, _ := strconv.ParseFloat(coords[0], 64)

		geoData := models.GeoJSON{
			Type: "FeatureCollection",
			Features: []models.GeoJSONFeature{
				models.GeoJSONFeature{
					Type: "Feature",
					Geometry: models.GeoJSONGeometry{
						Type:        "Point",
						Coordinates: []interface{}{lat, long},
					},
				},
			},
		}

		poi := models.Poi{
			Name:      record[0],
			Icon:      record[1],
			GeoData:   geoData,
			Active:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err := collections.GetPoicollection(*db.GetDbClient()).InsertOne(context.Background(), poi)
		if err != nil {
			panic(err)
		}
	}
}

func dataseed() {
	// create main admin
	passw, err := utils.HashPassword(config.GetEnv("DEFAULT_ADMIN_PASSW"))
	if err != nil {
		log.Fatal(err)
	}
	fabian := models.User{
		Email:     "fabian@fabianpoels.com",
		Password:  passw,
		Name:      "Fabian",
		Active:    true,
		Role:      "superadmin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	res, err := collections.GetUserCollection(*db.GetDbClient()).InsertOne(context.Background(), fabian)
	if err != nil {
		panic(err)
	}
	log.Println(res)

	// create categories
	// categories := []interface{}{
	// 	models.Category{
	// 		Key:    "red_zone",
	// 		Color:  "#f44336",
	// 		Icon:   "report",
	// 		Active: true,
	// 	},
	// 	models.Category{
	// 		Key:    "traffic_incident",
	// 		Color:  "#ff9100",
	// 		Icon:   "car_crash",
	// 		Active: true,
	// 	},

	// 	models.Category{
	// 		Key:    "protest",
	// 		Color:  "#cddc39",
	// 		Icon:   "groups",
	// 		Active: true,
	// 	},
	// 	models.Category{
	// 		Key:    "military",
	// 		Color:  "#cddc39",
	// 		Icon:   "radar",
	// 		Active: true,
	// 	},
	// 	models.Category{
	// 		Key:    "weather",
	// 		Color:  "#cddc39",
	// 		Icon:   "thunderstorm",
	// 		Active: true,
	// 	},
	// }

	// ress, err := models.GetCategoryCollection(*db.GetDbClient()).InsertMany(context.Background(), categories)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(ress)
}

func addUser() {
	passw, err := utils.HashPassword("")
	if err != nil {
		log.Fatal(err)
	}
	newUser := models.User{
		Email:     "mhannaralph@hotmail.com",
		Password:  passw,
		Name:      "Ralph",
		Active:    true,
		Role:      "superadmin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	res, err := collections.GetUserCollection(*db.GetDbClient()).InsertOne(context.Background(), newUser)
	if err != nil {
		panic(err)
	}
	log.Println(res)
}

type demoUser struct {
	Name     string
	Password string
}

func addDemoUsers() {
	demoUsers := []demoUser{
		{
			Name:     "demouser1",
			Password: "demouser1password",
		},
		{
			Name:     "demouser2",
			Password: "demouser2password",
		},
		{
			Name:     "demouser3",
			Password: "demouser3password",
		},
		{
			Name:     "demouser4",
			Password: "demouser4password",
		},
		{
			Name:     "demouser5",
			Password: "demouser5password",
		},
	}

	for _, user := range demoUsers {
		passw, err := utils.HashPassword(user.Password)
		if err != nil {
			log.Fatal(err)
		}
		newUser := models.User{
			Email:     fmt.Sprintf("%s@nonexistentmail.com", user.Name),
			Password:  passw,
			Name:      user.Name,
			Active:    true,
			Role:      "superadmin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		res, err := collections.GetUserCollection(*db.GetDbClient()).InsertOne(context.Background(), newUser)
		if err != nil {
			panic(err)
		}
		log.Println(res)
	}
}

func riskLevelsMigration() {
	mongoClient := db.GetDbClient()
	c := context.Background()

	// add risklevel to all municipalities
	update := bson.M{
		"$set": bson.M{
			"riskLevel": -1,
			"updatedAt": time.Now(),
		},
	}
	collections.GetMunicipalityCollection(*mongoClient).UpdateMany(c, bson.D{{}}, update, options.Update())

	// update if already present
	cur, err := collections.GetRiskLevelCollection(*mongoClient).Find(c, bson.D{{}}, options.Find())

	if err != nil {
		panic(err)
	}

	riskLevels := make([]models.RiskLevel, 0)

	err = cur.All(c, &riskLevels)
	if err != nil {
		panic(err)
	}

	for _, riskLevel := range riskLevels {
		log.Println(riskLevel.Level)
		municipalityObjectIDs := make([]primitive.ObjectID, len(riskLevel.Municipalities))

		for i, m := range riskLevel.Municipalities {
			municipalityObjectIDs[i] = m.Id
		}
		log.Println(municipalityObjectIDs)
		filter := bson.M{
			"_id": bson.M{"$in": municipalityObjectIDs},
		}
		update := bson.M{
			"$set": bson.M{
				"riskLevel": riskLevel.Level,
				"updatedAt": time.Now(),
			},
		}
		collections.GetMunicipalityCollection(*mongoClient).UpdateMany(c, filter, update, options.Update())
	}
}
