package main

import (
	. "github.com/yigitozgumus/leaderboard-api/server"
	"log"
	"net/http"
)

func main() {
	server := &LeaderboardServer{}
	log.Fatal(http.ListenAndServe(":5000", server))
}
