package datastore

import (
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
	"github.com/yigitozgumus/leaderboard-api/server"
	"math"
	"sync"
	"time"
)

//go:generate gotemplate "github.com/ncw/gotemplate/treemap" "RankingMap(float64, map[string]int)"

type InMemoryRankingLeaderboardStore struct {
	playerMap      map[string]server.User // user id -> user
	rankMap        *RankingMap            // score -> list of users with that score
	displayNameMap map[string]string      // user display name -> user id
	userLock       *sync.Mutex
	scoreLock      *sync.Mutex
}

func (i *InMemoryRankingLeaderboardStore) GetUserRankings() []server.User {
	var leaderboard []server.User
	for it := i.rankMap.Iterator(); it.Valid(); it.Next() {
		users := it.Value()
		for key := range users {
			user := i.playerMap[key]
			user.UserId = "" // remove userId on display
			user.Rank = uint32(len(leaderboard) + 1)
			leaderboard = append(leaderboard, user)
		}
	}
	return leaderboard
}

func (i *InMemoryRankingLeaderboardStore) GetUserRankingsFiltered(country string) []server.User {
	leaderboard := i.GetUserRankings()
	filtered := leaderboard[:0]
	for _, user := range leaderboard {
		if user.Country == country {
			filtered = append(filtered, user)
		}
	}
	return filtered
}

func (i *InMemoryRankingLeaderboardStore) CreateUserProfile(user server.User) error {
	user.UserId = uuid.New().String()
	if _, exists := i.displayNameMap[user.DisplayName]; exists {
		return server.UserExistsError
	}
	i.userLock.Lock()
	defer i.userLock.Unlock()
	i.displayNameMap[user.DisplayName] = user.UserId
	user.Rank = uint32(len(i.displayNameMap) + 1)
	i.playerMap[user.UserId] = user
	if rankings, exists := i.rankMap.Get(user.Points); exists {
		rankings[user.UserId] = 1
		i.rankMap.findNode(user.Points).value = rankings
	} else {
		i.rankMap.Set(user.Points, map[string]int{user.UserId: 1})
	}
	return nil
}

func (i *InMemoryRankingLeaderboardStore) GetUserProfile(userId string) (server.User, error) {
	user, exists := i.playerMap[userId]
	if exists == false {
		return server.User{}, server.NoUserPresentError
	}
	return user, nil
}

func (i *InMemoryRankingLeaderboardStore) SubmitUserScore(score server.Score) (server.Score, error) {
	score.TimeStamp = time.Now().String()
	if user, exists := i.playerMap[score.UserId]; exists == false {
		return score, server.NoUserPresentError
	} else {
		i.scoreLock.Lock()
		defer i.scoreLock.Unlock()
		currentScore := user.Points
		if users, exists := i.rankMap.Get(currentScore); exists {
			delete(users, user.UserId)
			i.rankMap.findNode(currentScore).value = users
		}
		newScore := currentScore + score.Score
		newScore = math.Round(newScore*100) / 100
		if users, exists := i.rankMap.Get(newScore); exists {
			users[user.UserId] = 1
			i.rankMap.findNode(newScore).value = users
		} else {
			i.rankMap.Set(newScore, map[string]int{user.UserId: 1})
		}
		user.Points = newScore
		i.playerMap[user.UserId] = user
	}
	return score, nil
}

func (i *InMemoryRankingLeaderboardStore) CreateUserProfiles(submission server.Submission) error {
	userSize := submission.SubmissionSize
	for index := 0; index < userSize; index++ {
		_ = i.CreateUserProfile(server.User{DisplayName: randstr.String(10), Country: getRandomEntry(countryList)})
	}
	return nil
}

func (i *InMemoryRankingLeaderboardStore) CreateScoreSubmissions(submission server.Submission) error {
	numberOfScores := submission.SubmissionSize
	userList := getUserList(i.playerMap)
	for index := 0; index < numberOfScores; index++ {
		score := getRandomScore(submission)
		_, _ = i.SubmitUserScore(server.Score{Score: score, UserId: getRandomEntry(userList)})
	}
	return nil
}

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
