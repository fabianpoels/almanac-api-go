package db

import (
	"almanac-api/config"
	"fmt"
	"log"

	"github.com/valkey-io/valkey-go"
)

var valkeyClient valkey.Client

func GetCacheClient() valkey.Client {
	if valkeyClient == nil {
		CacheConnect()
	}
	return valkeyClient
}

func CacheConnect() {
	host := config.GetEnv("VALKEY_HOST")
	port := config.GetEnv("VALKEY_PORT")
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{fmt.Sprintf("%s:%s", host, port)}})
	if err != nil {
		log.Fatal("⛒ Connection Failed to Cache")
		log.Fatal(err)
	}

	log.Println("Connected to cache")

	valkeyClient = client
}
