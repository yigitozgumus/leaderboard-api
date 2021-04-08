package server

import (
	"encoding/json"
	"errors"
	"net/http"
)

const jsonContentType = "application/json"
const keyContentType = "Content-Type"

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
	w.Header().Set(keyContentType, jsonContentType)
	if err := assertCorrectMethodType(r.Method, http.MethodGet); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(l.store.GetUserRankings())
}

// handles returning the current leaderboard filtered by the country (GET)
func (l *LeaderboardServer) leaderboardFilterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(keyContentType, jsonContentType)
	if err := assertCorrectMethodType(r.Method, http.MethodGet); err != nil {
		return
	}
	country := r.URL.Path[len("/leaderboard/"):]
	if len(country) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(invalidCountryError.Error())
		return
	}
	json.NewEncoder(w).Encode(l.store.GetUserRankingsFiltered(country))
}

// handles score submission of a user (POST)
func (l *LeaderboardServer) scoreSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	if err := assertCorrectMethodType(r.Method, http.MethodPost); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	successResponse(w)
}

// handles returning the user profile with given guid (GET)
func (l *LeaderboardServer) userProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(keyContentType, jsonContentType)
	if err := assertCorrectMethodType(r.Method, http.MethodGet); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// handles creating user with given information (POST)
func (l *LeaderboardServer) createUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := assertCorrectMethodType(r.Method, http.MethodPost); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	var u User
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&u)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	l.store.CreateUserProfile(u)
	successResponse(w)
}

func assertCorrectMethodType(requestType string, methodType string) error {
	if requestType != methodType {
		return invalidRequestTypeError
	}
	return nil
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func successResponse(w http.ResponseWriter) {
	errorResponse(w, "Success", http.StatusOK)
}
