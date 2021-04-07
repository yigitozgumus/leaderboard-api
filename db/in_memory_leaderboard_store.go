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
	var leaderboard []User
	for _, user := range i.store {
		leaderboard = append(leaderboard, User{DisplayName: user.DisplayName, Points: user.Points, Rank: user.Rank, Country: user.Country})
	}
	return leaderboard
}

// FIXME add filter for country
func (i *InMemoryLeaderboardStore) GetUserRankingsFiltered(country string) []User {
	var leaderboard []User
	for _, user := range i.store {
		leaderboard = append(leaderboard, User{DisplayName: user.DisplayName, Points: user.Points, Rank: user.Rank, Country: user.Country})
	}
	return leaderboard
}

func (i *InMemoryLeaderboardStore) CreateUserProfile(user User) {
	i.store = append(i.store, user)
}

func (i *InMemoryLeaderboardStore) GetUserProfile(name string) User {

	for _, user := range i.store {
		if user.DisplayName == name {
			return user
		}
	}
	// return empty user
	return User{}
}

func NewInMemoryLeaderboardStore() *InMemoryLeaderboardStore {
	return &InMemoryLeaderboardStore{
		nil,
		&sync.Mutex{},
	}
}
