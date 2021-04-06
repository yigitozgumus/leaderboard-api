package server

import "net/http"

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
	w.WriteHeader(http.StatusOK)
}