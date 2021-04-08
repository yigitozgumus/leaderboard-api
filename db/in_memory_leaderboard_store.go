package db

import (
	. "github.com/yigitozgumus/leaderboard-api/server"
	"sync"
)

type InMemoryLeaderboardStore struct {
	store []User
	mu    *sync.Mutex
}

func (i *InMemoryLeaderboardStore) GetUserRankings() []User {
	return i.store
}

func (i *InMemoryLeaderboardStore) GetUserRankingsFiltered(country string) []User {
	leaderboard := i.store[:0]
	for _, user := range i.store {
		if user.Country == country {
			leaderboard = append(leaderboard, user)
		}
	}
	return leaderboard
}

func (i *InMemoryLeaderboardStore) CreateUserProfile(user User) {
	// FIXME handle ranking
	i.store = append(i.store, user)
}

// FIXME add error for no user present
func (i *InMemoryLeaderboardStore) GetUserProfile(name string) (User, error) {

	for _, user := range i.store {
		if user.DisplayName == name {
			return user, nil
		}
	}
	// return empty user
	return User{}, nil
}

func NewInMemoryLeaderboardStore() *InMemoryLeaderboardStore {
	return &InMemoryLeaderboardStore{
		nil,
		&sync.Mutex{},
	}
}
