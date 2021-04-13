package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/yigitozgumus/leaderboard-api/datastore"
	"github.com/yigitozgumus/leaderboard-api/server"
	"log"
	"net/http"
	"os"
)

const serverTypeDev = "dev"

func main() {
	// load .env file
	_ = godotenv.Load(".env")

	serverType := os.Getenv("SERVER")
	storageType := os.Getenv("STORAGE_TYPE")
	var redisURL string
	if serverType == "dev" {
		redisURL = os.Getenv("REDIS_URL")
	} else {
		parse, _ := redis.ParseURL(os.Getenv("REDIS_URL"))
		redisURL = parse.Addr
	}
	//redisUrl, _ := redis.ParseURL(os.Getenv("REDIS_URL"))
	configuration := server.ConfigurationType{
		Server:     serverType,
		Storage:    storageType,
		MongoUri: os.Getenv("ATLAS_URL"),
		RedisUri: redisURL,
		Message:    "Initializing " + serverType + " server with " + storageType,
	}
	var leaderboardServer *server.LeaderboardServer
	switch configuration.Storage {
	case "memory":
		store := datastore.NewInMemoryRankingStore()
		leaderboardServer = ConfigureServer(configuration, store)
	case "final":
		store := datastore.NewDatabaseLeaderboardStore(configuration)
		leaderboardServer = ConfigureServer(configuration, store)
		closeConnection := store.InitializeConnection()
		store.InitializeRedisCache()
		defer closeConnection()
	}
	var port = os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatal(http.ListenAndServe(":" + port, leaderboardServer))
}


func ConfigureServer(s server.ConfigurationType, store server.LeaderboardStore) *server.LeaderboardServer {
	var leaderboardServer *server.LeaderboardServer
	if s.Server == serverTypeDev {
		leaderboardServer = server.NewLeaderboardServer(store, true)
	} else {
		leaderboardServer = server.NewLeaderboardServer(store, false)
	}
	fmt.Println(s.Message)
	return leaderboardServer
}
