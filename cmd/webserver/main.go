package main

import (
	"log"
	"net/http"

	"github.com/yigitozgumus/leaderboard-api/db"
	"github.com/yigitozgumus/leaderboard-api/server"
)

func main() {
	store := db.NewInMemoryRankingStore()
	leaderboardServer := server.NewLeaderboardServer(store)
	log.Fatal(http.ListenAndServe(":5000", leaderboardServer))
}
