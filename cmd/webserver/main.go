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
		leaderboardServer = server.NewLeaderBoardProductionServer(store)
		fmt.Println("Initializing Production Server")
	} else {
		leaderboardServer = server.NewLeaderBoardDevelopmentServer(store)
		fmt.Println("Initializing Development Server")
	}

	log.Fatal(http.ListenAndServe(":5000", leaderboardServer))
}
