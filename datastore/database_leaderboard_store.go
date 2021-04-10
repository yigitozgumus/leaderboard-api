package datastore

import (
	"context"
	"github.com/yigitozgumus/leaderboard-api/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type DatabaseLeaderboardStore struct {
	Client *mongo.Client
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

func NewDatabaseLeaderboardStore(config server.ConfigurationType) (*DatabaseLeaderboardStore, func()) {
	clientOptions := options.Client().ApplyURI(config.Connection)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	closeConnection := func() {
		client.Disconnect(context.TODO())
	}
	return &DatabaseLeaderboardStore{client}, closeConnection
}
