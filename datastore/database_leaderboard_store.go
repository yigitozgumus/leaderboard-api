package datastore

import (
	"github.com/yigitozgumus/leaderboard-api/server"
	"io"
)

type DatabaseLeaderboardStore struct {
	playerDatabase  io.Reader
	rankingDatabase io.Reader
}

func (f *DatabaseLeaderboardStore) GetUserRankings() []server.User {
	// FIXME
	return nil
}

func (f *DatabaseLeaderboardStore) GetUserRankingsFiltered(country string) []server.User {
	// FIXME
	return nil
}
func (f *DatabaseLeaderboardStore) CreateUserProfile(user server.User) error {
	// FIXME
	return nil
}
func (f *DatabaseLeaderboardStore) GetUserProfile(userId string) (server.User, error) {
	// FIXME
	return server.User{}, nil
}

func (f *DatabaseLeaderboardStore) SubmitUserScore(score server.Score) (server.Score, error) {
	// FIXME
	return server.Score{}, nil
}
func (f *DatabaseLeaderboardStore) CreateUserProfiles(submission server.Submission) error {
	// FIXME
	return nil
}
func (f *DatabaseLeaderboardStore) CreateScoreSubmissions(submission server.Submission) error {
	// FIXME
	return nil
}
