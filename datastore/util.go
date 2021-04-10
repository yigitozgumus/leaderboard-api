package datastore

import (
	"github.com/yigitozgumus/leaderboard-api/server"
	"math/rand"
	"time"
)

// helpers
var countryList = []string{"tr", "de", "fr", "au", "us"}

func scoreRankCompare(x, y float64) bool { return x > y }

func getRandomEntry(list []string) string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })
	return list[0]
}

func getUserList(players map[string]server.User) []string {
	var leaderboard []string
	for key := range players {
		leaderboard = append(leaderboard, key)
	}
	return leaderboard
}

func getRandomScore(s server.Submission) float64 {
	return float64(rand.Intn(s.MaxScore-s.MinScore) + s.MinScore)
}
