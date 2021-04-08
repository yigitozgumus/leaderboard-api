package db

import (
	"github.com/google/uuid"
	"github.com/yigitozgumus/leaderboard-api/server"
	"sync"
)

type InMemoryLeaderboardStore struct {
	store []server.User
	mu    *sync.Mutex
}

func (i *InMemoryLeaderboardStore) GetUserRankings() []server.User {
	return i.store
}

func (i *InMemoryLeaderboardStore) GetUserRankingsFiltered(country string) []server.User {
	leaderboard := i.store[:0]
	for _, user := range i.store {
		if user.Country == country {
			leaderboard = append(leaderboard, user)
		}
	}
	return leaderboard
}

func (i *InMemoryLeaderboardStore) CreateUserProfile(user server.User) error {
	for _, u := range i.store {
		if u.DisplayName == user.DisplayName {
			return server.UserExistsError
		}
	}
	user.UserId = uuid.New().String()
	i.store = append(i.store, user)
	return nil
}

// FIXME add error for no user present
func (i *InMemoryLeaderboardStore) GetUserProfile(name string) (server.User, error) {

	for _, user := range i.store {
		if user.UserId == name {
			return user, nil
		}
	}
	return server.User{}, server.NoUserPresentError
}

func NewInMemoryLeaderboardStore() *InMemoryLeaderboardStore {
	return &InMemoryLeaderboardStore{
		nil,
		&sync.Mutex{},
	}
}
