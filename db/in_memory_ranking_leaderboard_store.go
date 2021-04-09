package db

import (
	"github.com/google/uuid"
	"github.com/yigitozgumus/leaderboard-api/server"
	"sync"
)

//go:generate gotemplate "github.com/ncw/gotemplate/treemap" "RankingMap(float64, []string)"

type InMemoryRankingLeaderboardStore struct {
	playerMapStore map[string]server.User // user id -> user
	rankMapStore *RankingMap // score -> list of users with that score
	idDisplayMap map[string]string // user display name -> user id
	userLock *sync.Mutex
	scoreLock *sync.Mutex
}

func (i *InMemoryRankingLeaderboardStore) GetUserRankings() []server.User {
	var leaderboard []server.User
	for it := i.rankMapStore.Iterator(); it.Valid(); it.Next() {
		users := it.Value()
		for _, u := range users {
			leaderboard = append(leaderboard, i.playerMapStore[u])
		}
	}
	return leaderboard
}

func (i *InMemoryRankingLeaderboardStore) GetUserRankingsFiltered(country string) []server.User {
	panic("implement me")
}

func (i *InMemoryRankingLeaderboardStore) CreateUserProfile(user server.User) error {
	user.UserId = uuid.New().String()
	if _, exists := i.idDisplayMap[user.DisplayName]; exists {
		return server.UserExistsError
	}
	i.userLock.Lock()
	defer i.userLock.Unlock()
	i.idDisplayMap[user.DisplayName] = user.UserId
	i.playerMapStore[user.UserId] = user
	if rankings, exists := i.rankMapStore.Get(user.Points); exists {
		rankings = append(rankings, user.UserId)
		i.rankMapStore.findNode(user.Points).value = rankings
	} else {
		i.rankMapStore.Set(user.Points, []string{user.UserId})
	}
	return nil
}

func (i *InMemoryRankingLeaderboardStore) GetUserProfile(userId string) (server.User, error) {
	user, exists := i.playerMapStore[userId]
	if exists == false {
		return server.User{}, server.NoUserPresentError
	}
	return user, nil
}

func (i *InMemoryRankingLeaderboardStore) SubmitUserScore(score server.Score) (server.Score, error) {
	panic("implement me")
}

func scoreRankCompare(x, y float64) bool { return x > y }

func NewInMemoryRankingStore() *InMemoryRankingLeaderboardStore {
	rankMap := NewRankingMap(scoreRankCompare)
	return &InMemoryRankingLeaderboardStore{
		map[string]server.User{},
		rankMap,
		map[string]string{},
		&sync.Mutex{},
		&sync.Mutex{},
	}
}