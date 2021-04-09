package db

import (
	"github.com/yigitozgumus/leaderboard-api/server"
	"sync"
)

type InMemoryRankingLeaderboardStore struct {
	playerMapStore map[string]server.User
	rankMapStore map[float64][]string
	userLock *sync.Mutex
	scoreLock *sync.Mutex
}

func (i *InMemoryRankingLeaderboardStore) GetUserRankings() []server.User {
	panic("implement me")
}

func (i *InMemoryRankingLeaderboardStore) GetUserRankingsFiltered(country string) []server.User {
	panic("implement me")
}

func (i *InMemoryRankingLeaderboardStore) CreateUserProfile(user server.User) error {
	panic("implement me")
}

func (i *InMemoryRankingLeaderboardStore) GetUserProfile(name string) (server.User, error) {
	panic("implement me")
}

func (i *InMemoryRankingLeaderboardStore) SubmitUserScore(score server.Score) (server.Score, error) {
	panic("implement me")
}

func NewInMemoryRankingStore() *InMemoryRankingLeaderboardStore {
	return &InMemoryRankingLeaderboardStore{
		map[string]server.User{},
		map[float64][]string{},
		&sync.Mutex{},
		&sync.Mutex{},
	}
}