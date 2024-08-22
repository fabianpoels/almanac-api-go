package main

import (
	"flag"
	"fmt"
	"os"

	"almanac-api/config"
	"almanac-api/db"
	"almanac-api/server"
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

// func dataseed() {
// 	// create main admin
// 	passw, err := utils.HashPassword("Test123Test123")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fabian := models.User{
// 		Email:     "fabian@fabianpoels.com",
// 		Password:  passw,
// 		Name:      "Fabian",
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}
// 	res, err := models.GetUserCollection(*db.GetDbClient()).InsertOne(context.Background(), fabian)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Println(res)
// }
