package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"almanac-api/collections"
	"almanac-api/config"
	"almanac-api/db"
	"almanac-api/server"
	"almanac-api/utils"

	"gitlab.com/almanac-app/models"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	db.DbConnect()
	db.CacheConnect()

	// dataseeding
	// dataseed()

	// start server
	server.Init()
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
