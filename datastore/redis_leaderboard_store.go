package datastore

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
	"github.com/yigitozgumus/leaderboard-api/server"
	"golang.org/x/net/context"
	"sync"
	"time"
)

var (
	ErrNil = errors.New("no matching record found in redis database")
	Ctx    = context.TODO()
)

var leaderboardKey = "leaderboard"
var userKey = "user:"

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
	if err := client.Ping(Ctx).Err(); err != nil {
		return err
	}
	r.Client = client
	return nil
}

func (r *RedisLeaderboardStore) GetUserRankings() []server.User {
	scores := r.Client.ZRevRangeWithScores(Ctx, leaderboardKey, 0, -1)
	if scores == nil {
		return nil
	}
	count := len(scores.Val())
	users := make([]server.User, count)
	for idx, member := range scores.Val() {
		users[idx] = server.User{
			DisplayName: member.Member.(string),
			Points: float64(member.Score),
			Rank: int64(idx +1),
		}
	}
	return users
}

func (r *RedisLeaderboardStore) GetUserRankingsFiltered(country string) []server.User {
	// FIXME
	return nil
}

func (r *RedisLeaderboardStore) CreateUserProfile(user server.User) error {
	user.UserId = uuid.New().String()
	pipe := r.Client.TxPipeline()
	pipe.Get(Ctx,user.DisplayName)
	_, err := pipe.Exec(Ctx)
	if err == nil {
		return server.UserExistsError
	}
	u, _ := user.MarshalBinary()
	member := &redis.Z{
		Score: float64(user.Points),
		Member: user.DisplayName,
	}
	pipe.Set(Ctx, user.DisplayName, userKey+user.UserId , -1)
	pipe.Set(Ctx, userKey+user.UserId, u, -1)
	pipe.ZAdd(Ctx, leaderboardKey, member)
	pipe.ZRank(Ctx, leaderboardKey, user.DisplayName)
	_, err = pipe.Exec(Ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisLeaderboardStore) GetUserProfile(userId string) (server.User, error) {
	pipe := r.Client.TxPipeline()
	res := pipe.Get(Ctx,userKey+userId)
	_, err := pipe.Exec(Ctx)
	if err != nil {
		return server.User{}, err
	}
	var u server.User
	user,_ := res.Bytes()
	json.Unmarshal(user, &u)
	rank := pipe.ZRevRank(Ctx, leaderboardKey, u.DisplayName)
	score := pipe.ZScore(Ctx, leaderboardKey, u.DisplayName)
	_, err = pipe.Exec(Ctx)
	if err != nil {
		return server.User{}, err
	}
	u.Rank, _ = rank.Result()
	u.Rank +=1
	u.Points, _ = score.Result()
	return u, nil
}

func (r *RedisLeaderboardStore) SubmitUserScore(score server.Score) (server.Score, error) {
	score.TimeStamp = time.Now().String()
	pipe := r.Client.TxPipeline()
	res := pipe.Get(Ctx,userKey+score.UserId)
	_, err := pipe.Exec(Ctx)
	if err != nil {
		return server.Score{}, err
	}
	var u server.User
	user, _ := res.Bytes()
	json.Unmarshal(user, &u)
	u.Points += score.Score
	pipe.Set(Ctx, userKey+u.UserId, u, -1)
	member := &redis.Z{
		Score: float64(u.Points),
		Member: u.DisplayName,
	}
	pipe.ZAdd(Ctx, leaderboardKey, member)
	_, err = pipe.Exec(Ctx)
	if err != nil {
		return server.Score{}, err
	}
	return score, nil
}

func (r *RedisLeaderboardStore) CreateUserProfiles(submission server.Submission) error {
	userSize := submission.SubmissionSize
	for index := 0; index < userSize; index++ {
		_ = r.CreateUserProfile(server.User{DisplayName: randstr.String(10), Country: getRandomEntry(countryList)})
	}
	return nil
}

func (r *RedisLeaderboardStore) CreateScoreSubmissions(submission server.Submission) error {
	users := r.Client.Keys(Ctx, userKey+"*")
	userList := make([]string, len(users.Val()))
	for ind, user := range users.Val() {
		userList[ind] = user
	}
	numberOfScores := submission.SubmissionSize
	for index := 0; index < numberOfScores; index++ {
		score := getRandomScore(submission)
		_, _ = r.SubmitUserScore(server.Score{Score: score, UserId: getRandomEntry(userList)[5:]})
	}
	return nil
}