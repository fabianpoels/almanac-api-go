package db

import (
	"almanac-api/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func GetCacheClient() redis.Client {
	if redisClient == nil {
		CacheConnect()
	}
	return *redisClient
}

func CacheConnect() {
	host := config.GetEnv("REDIS_HOST")
	port := config.GetEnv("REDIS_PORT")

	log.Printf("Connecting to redis client at: %s", fmt.Sprintf("%s:%s", host, port))
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println(err)
		log.Fatal("⛒ Connection Failed to Cache")
		log.Fatal(err)
	}
	defer cancel()

	log.Println("Connected to cache")

	redisClient = client
}
