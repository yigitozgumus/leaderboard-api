package datastore

import (
	"errors"
	"github.com/yigitozgumus/leaderboard-api/server"
	"golang.org/x/net/context"
	"sync"
)
import "github.com/go-redis/redis"

var (
	ErrNil = errors.New("no matching record found in redis database")
	Ctx    = context.TODO()
)

type RedisLeaderboardStore struct {
	Client     *redis.Client
	Connection string
	userLock   *sync.Mutex
	scoreLock  *sync.Mutex
}

func NewRedisLeaderboardStore(config server.ConfigurationType) *RedisLeaderboardStore {
	return &RedisLeaderboardStore{
		nil,
		config.Connection,
		&sync.Mutex{},
		&sync.Mutex{},
	}
}

func (r *RedisLeaderboardStore) InitializeConnection() error {
	client := redis.NewClient(&redis.Options{
		Addr: r.Connection,
		Password: "",
		DB: 0,
	})
	if err := client.Ping().Err(); err != nil {
		return err
	}
	r.Client = client
	return nil
}

func (r *RedisLeaderboardStore) GetUserRankings() []server.User {
	// FIXME
	return nil
}

func (r *RedisLeaderboardStore) GetUserRankingsFiltered(country string) []server.User {
	// FIXME
	return nil
}

func (r *RedisLeaderboardStore) CreateUserProfile(user server.User) error {
	// FIXME
	return nil
}

func (r *RedisLeaderboardStore) GetUserProfile(userId string) (server.User, error) {
	// FIXME
	return server.User{},nil
}

func (r *RedisLeaderboardStore) SubmitUserScore(score server.Score) (server.Score, error) {
	// FIXME
	return server.Score{}, nil
}

func (r *RedisLeaderboardStore) CreateUserProfiles(submission server.Submission) error {
	// FIXME
	return nil
}

func (r *RedisLeaderboardStore) CreateScoreSubmissions(submission server.Submission) error {
	// FIXME
	return nil
}