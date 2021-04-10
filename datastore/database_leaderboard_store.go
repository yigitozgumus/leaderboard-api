package datastore

import (
	"github.com/yigitozgumus/leaderboard-api/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseLeaderboardStore struct {

}

func (d *DatabaseLeaderboardStore) configureDatabaseConnection() {
	// FIXME
}

func (d *DatabaseLeaderboardStore) GetUserRankings() []server.User {
	// FIXME
	return nil
}

func (d *DatabaseLeaderboardStore) GetUserRankingsFiltered(country string) []server.User {
	// FIXME
	return nil
}
func (d *DatabaseLeaderboardStore) CreateUserProfile(user server.User) error {
	// FIXME
	return nil
}
func (d *DatabaseLeaderboardStore) GetUserProfile(userId string) (server.User, error) {
	// FIXME
	return server.User{}, nil
}

func (d *DatabaseLeaderboardStore) SubmitUserScore(score server.Score) (server.Score, error) {
	// FIXME
	return server.Score{}, nil
}
func (d *DatabaseLeaderboardStore) CreateUserProfiles(submission server.Submission) error {
	// FIXME
	return nil
}
func (d *DatabaseLeaderboardStore) CreateScoreSubmissions(submission server.Submission) error {
	// FIXME
	return nil
}

func NewDatabaseLeaderboardStore() (*DatabaseLeaderboardStore, func()) {
	//FIXME
}
