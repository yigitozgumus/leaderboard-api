package server

import "net/http"

func LeaderboardServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}