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
	configuration := server.ConfigurationType{
		Server:     os.Getenv("SERVER"),
		Connection: os.Getenv("URI"),
		Message:    "Initializing Development Server",
	}

	store := datastore.NewDatabaseLeaderboardStore(configuration)
	closeConnection := store.InitializeConnection()
	defer closeConnection()
	var leaderboardServer *server.LeaderboardServer
	leaderboardServer = ConfigureServer(configuration, store)
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
