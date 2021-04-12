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
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("Error loading .env file")
	}
	serverType := os.Getenv("SERVER")
	storageType := os.Getenv("STORAGE_TYPE")
	redisUrl, _ := redis.ParseURL(os.Getenv("REDIS_URL"))
	configuration := server.ConfigurationType{
		Server:     serverType,
		Storage:    storageType,
		MongoUri: os.Getenv("ATLAS_URL"),
		RedisUri: redisUrl.Addr,
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

	log.Fatal(http.ListenAndServe(":5000", leaderboardServer))
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
