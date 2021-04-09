package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/yigitozgumus/leaderboard-api/db"
	"github.com/yigitozgumus/leaderboard-api/server"
)

func main() {
	store := db.NewInMemoryRankingStore()
	var leaderboardServer *server.LeaderboardServer
	argDoesNotExist := len(os.Args) == 1
	if argDoesNotExist{
		leaderboardServer = server.NewLeaderboardServer(store, false)
		fmt.Println("Initializing Production Server")
	} else if arg := os.Args[1]; arg == "dev" {
		leaderboardServer = server.NewLeaderboardServer(store, true)
		fmt.Println("Initializing Development Server")
	} else {
		fmt.Println("Wrong Parameter, Did you mean \"dev\" ?")
		os.Exit(1)
	}

	log.Fatal(http.ListenAndServe(":5000", leaderboardServer))
}
