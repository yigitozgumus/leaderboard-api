package main

import (
	S "github.com/yigitozgumus/leaderboard-api/server"
	"log"
	"net/http"
)

func main() {

	handler := http.HandlerFunc(S.LeaderboardServer)
	log.Fatal(http.ListenAndServe(":5000", handler))
}
