package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
)

// consts
const jsonContentType = "application/json"
const keyContentType = "Content-Type"

type User struct {
	UserId      string  `json:"user_id"`
	DisplayName string  `json:"display_name"`
	Points      float64 `json:"points"`
	Rank        uint32  `json:"rank"`
	Country     string  `json:"country"`
}

type Score struct {
	Score     float64 `json:"score"`
	UserId    string  `json:"user_id"`
	TimeStamp string  `json:"time_stamp"`
}

type LeaderboardStore interface {
	GetUserRankings() []User
	GetUserRankingsFiltered(country string) []User
	CreateUserProfile(user User) error
	GetUserProfile(userId string) (User, error)
	SubmitUserScore(score Score) (Score, error)
	CreateUserProfiles(submission Submission) error
	CreateScoreSubmissions(submission Submission) error
}

type LeaderboardServer struct {
	store LeaderboardStore
	http.Handler
}

type Submission struct {
	SubmissionSize int  `json:"submission_size"`
	MaxScore       int `json:"max_score"`
	MinScore       int `json:"min_score"`
}

// errors
var invalidCountryError = errors.New("invalid country input")
var UserExistsError = errors.New("user exists")
var NoUserPresentError = errors.New("no user present")

func NewLeaderboardServer(store LeaderboardStore, isDevelopment bool) *LeaderboardServer {
	l := new(LeaderboardServer)
	l.store = store
	router := chi.NewRouter()

	router.Get("/leaderboard", l.leaderboardHandler)
	router.Get("/leaderboard/{slug}", l.leaderboardFilterHandler)
	router.Post("/leaderboard/score/submit", l.scoreSubmissionHandler)
	router.Get("/leaderboard/user/profile/{slug}", l.userProfileHandler)
	router.Post("/leaderboard/user/create", l.createUserHandler)
	if isDevelopment {
		router.Post("/leaderboard/test/create-users", l.dummyUserHandler)
		router.Post("/leaderboard/test/submit-scores", l.dummyScoreSubmissionHandler)
	}

	l.Handler = router

	return l
}

// Handles returning the current leaderboard (GET)
func (l *LeaderboardServer) leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(keyContentType, jsonContentType)
	err := json.NewEncoder(w).Encode(l.store.GetUserRankings())
	if err != nil {
		panic(err)
	}
}

// handles returning the current leaderboard filtered by the country (GET)
func (l *LeaderboardServer) leaderboardFilterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(keyContentType, jsonContentType)
	country := chi.URLParam(r, "slug")
	if len(country) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(invalidCountryError.Error())
		if err != nil {
			panic(err)
		}
		return
	}
	err := json.NewEncoder(w).Encode(l.store.GetUserRankingsFiltered(country))
	if err != nil {
		panic(err)
	}
}

// handles score submission of a user (POST)
func (l *LeaderboardServer) scoreSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	var s Score
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&s)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	var score Score
	score, err = l.store.SubmitUserScore(s)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	scoreSubmissionResponse(w, score)
}

// handles returning the user profile with given guid (GET)
func (l *LeaderboardServer) userProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(keyContentType, jsonContentType)
	guid := chi.URLParam(r, "slug")
	user, err := l.store.GetUserProfile(guid)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		panic(err)
	}
}

// handles creating user with given information (POST)
func (l *LeaderboardServer) createUserHandler(w http.ResponseWriter, r *http.Request) {
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
	err = l.store.CreateUserProfile(u)
	if err != nil {
		if errors.As(err, &UserExistsError) {
			errorResponse(w, err.Error(), http.StatusForbidden)
		}
		return
	}
	successResponse(w)
}

// handles creating dummy users for testing
func (l *LeaderboardServer) dummyUserHandler(w http.ResponseWriter, r *http.Request) {
	var s Submission
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&s)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	err = l.store.CreateUserProfiles(s)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// handles creating dummy user score submissions for testing
func (l *LeaderboardServer) dummyScoreSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	var s Submission
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&s)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}
	err = l.store.CreateScoreSubmissions(s)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	_, err := w.Write(jsonResp)
	if err != nil {
		panic(err)
	}
}

func successResponse(w http.ResponseWriter) {
	errorResponse(w, "Success", http.StatusOK)
}

func scoreSubmissionResponse(w http.ResponseWriter, score Score) {
	w.Header().Set("Content-Type", jsonContentType)

	resp := make(map[string]Score)
	resp["submission"] = score
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusForbidden)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		panic(err)
	}
}
