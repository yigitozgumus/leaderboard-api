package server

import (
	"encoding/json"
	"errors"
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

// errors
var invalidCountryError = errors.New("invalid country input")
var invalidRequestTypeError = errors.New("invalid request type")

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

// Handles returning the current leaderboard (GET)
func (l *LeaderboardServer) leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	if err := assertCorrectMethodType(w, r.Method, http.MethodGet); err != nil {
		return
	}
	json.NewEncoder(w).Encode(l.store.GetUserRankings())
	w.WriteHeader(http.StatusOK)
}

// handles returning the current leaderboard filtered by the country (GET)
func (l *LeaderboardServer) leaderboardFilterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	if err := assertCorrectMethodType(w, r.Method, http.MethodGet); err != nil {
		return
	}
	country := r.URL.Path[len("/leaderboard/"):]
	if len(country) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(invalidCountryError.Error())
		return
	}
	json.NewEncoder(w).Encode(l.store.GetUserRankingsFiltered(country))
	w.WriteHeader(http.StatusOK)
}

// handles score submission of a user (POST)
func (l *LeaderboardServer) scoreSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	if err := assertCorrectMethodType(w, r.Method, http.MethodPost); err != nil {

	}
}

// handles returning the user profile with given guid (GET)
func (l *LeaderboardServer) userProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	if err := assertCorrectMethodType(w, r.Method, http.MethodGet); err != nil {
		return
	}
}

// handles creating user with given information (POST)
func (l *LeaderboardServer) createUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := assertCorrectMethodType(w, r.Method, http.MethodPost); err != nil {
		return
	}
}

func assertCorrectMethodType(w http.ResponseWriter, requestType string, methodType string) error {
	if requestType != methodType {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(invalidRequestTypeError.Error())
		return invalidRequestTypeError
	}
	return nil
}
