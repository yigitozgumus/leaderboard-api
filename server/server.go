package server

import (
	"net/http"
)

const prefix = "/leaderboard"

type User struct {
	displayName string
	points float64
	rank uint32
	country string
}

type UserStore interface {
	getUserRankings()
}

type LeaderboardServer struct {
	store UserStore
}

func (l *LeaderboardServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != prefix {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}