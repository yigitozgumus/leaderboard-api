package db

import (
	"github.com/google/uuid"
	"github.com/yigitozgumus/leaderboard-api/server"
	"sync"
	"time"
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
	i.mu.Lock()
	defer i.mu.Unlock()
	i.store = append(i.store, user)
	return nil
}

func (i *InMemoryLeaderboardStore) GetUserProfile(name string) (server.User, error) {

	for _, user := range i.store {
		if user.UserId == name {
			return user, nil
		}
	}
	return server.User{}, server.NoUserPresentError
}

func (i *InMemoryLeaderboardStore) SubmitUserScore(score server.Score) (server.Score, error) {
	score.TimeStamp = time.Now().String()
	index := -1
	for i, user := range i.store {
		if user.UserId == score.UserId {
			index = i
		}
	}
	if index != -1  {
		i.mu.Lock()
		defer i.mu.Unlock()
		i.store[index].Points = score.Score
		return score, nil
	}
	return score, server.NoUserPresentError
}

func NewInMemoryLeaderboardStore() *InMemoryLeaderboardStore {
	return &InMemoryLeaderboardStore{
		nil,
		&sync.Mutex{},
	}
}
