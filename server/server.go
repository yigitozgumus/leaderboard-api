package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const jsonContentType = "application/json"

type User struct {
	DisplayName string
	Points      float64
	Rank        uint32
	Country     string
}

type LeaderboardStore interface {
	GetUserRankings() []User
	GetUserRankingsFiltered(country string) []User
	CreateUserProfile(user User)
	GetUserProfile(name string) User // FIXME
}

type LeaderboardServer struct {
	store LeaderboardStore
	http.Handler
}

func NewLeaderboardServer(store LeaderboardStore) *LeaderboardServer {
	l := new(LeaderboardServer)
	l.store = store
	router := http.NewServeMux()

	router.Handle("/leaderboard", http.HandlerFunc(l.leaderboardHandler))
	router.Handle("/leaderboard/", http.HandlerFunc(l.leaderboardFilterHandler))
	router.Handle("/leaderboard/score/submit", http.HandlerFunc(l.scoreSubmissionHandler))
	router.Handle("/leaderboard/user/profile/", http.HandlerFunc(l.userProfileHandler))
	router.Handle("/leaderboard/user/create", http.HandlerFunc(l.createUserHandler))
	l.Handler = router

	return l
}

func (l *LeaderboardServer) leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(l.store.GetUserRankings())
	w.WriteHeader(http.StatusOK)
}

func (l *LeaderboardServer) leaderboardFilterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	country := r.URL.Path[len("/leaderboard/"):]
	if len(country) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "country code must not be null")
		return
	}
	json.NewEncoder(w).Encode(l.store.GetUserRankingsFiltered(country))
	w.WriteHeader(http.StatusOK)
}

func (l *LeaderboardServer) scoreSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	// FIXME
}

func (l *LeaderboardServer) userProfileHandler(w http.ResponseWriter, r *http.Request) {
	// FIXME
}

func (l *LeaderboardServer) createUserHandler(w http.ResponseWriter, r *http.Request) {

}
