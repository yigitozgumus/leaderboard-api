package main

import (
	"fmt"
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
		log.Fatalf("Error loading .env file")
	}
	serverType := os.Getenv("SERVER")
	storageType := os.Getenv("STORAGE_TYPE")
	configuration := server.ConfigurationType{
		Server:     serverType,
		Storage:    storageType,
		Connection: os.Getenv("URI"),
		Message:    "Initializing " + serverType + " server with " + storageType,
	}
	var leaderboardServer *server.LeaderboardServer
	switch configuration.Storage {
	case "memory":
		store := datastore.NewInMemoryRankingStore()
		leaderboardServer = ConfigureServer(configuration, store)
	case "mongo":
		store := datastore.NewDatabaseLeaderboardStore(configuration)
		leaderboardServer = ConfigureServer(configuration, store)
		closeConnection := store.InitializeConnection()
		defer closeConnection()

	case "redis":
		store := datastore.NewRedisLeaderboardStore(configuration)
		leaderboardServer = ConfigureServer(configuration, store)
		err := store.InitializeConnection()
		if err != nil {
			log.Fatalf("Failed to connect to redis: %s", err.Error())
		}
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
