package datastore

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
	"github.com/yigitozgumus/leaderboard-api/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	ErrNil = errors.New("no matching record found in redis database")
	Ctx    = context.TODO()
	dbError = errors.New("error with database")
	ErrScoreSubmission = errors.New("Score submission failure")
	ErrRankAcquire = errors.New("rank of the user is not available")
	leaderboardKey = "leaderboard"
	userKey = "user:"
)

type DatabaseLeaderboardStore struct {
	MongoClient *mongo.Client
	RedisClient *redis.Client
	mongoUri    string
	redisUri string
	userLock    *sync.Mutex
	scoreLock   *sync.Mutex
}

type Ranking struct {
	Score float64 `bson:"score"`
	CurrentUsers []UserRank `bson:"current_users"`
}

type UserRank struct {
	User string `bson:"user_id"`
}

func (d *DatabaseLeaderboardStore) GetUserRankings() []server.User {
	scores, err := d.RedisClient.ZRevRangeWithScores(Ctx, leaderboardKey, 0, -1).Result()
	if err != nil {
		log.Fatal(err)
	}
	var users []server.User
	for _, score := range scores {
		str, _ := score.Member.(string)
		var user server.User
		findOptions := options.FindOne()
		err = d.getUsers().FindOne(Ctx, bson.M{"display_name": str}, findOptions).Decode(&user)
		rank, err := d.RedisClient.ZRevRank(Ctx, leaderboardKey, user.DisplayName).Result()
		if err != nil {
			return nil
		}
		user.Rank = int64(rank + 1)
		users = append(users, user)
	}
	return users
}

func (d *DatabaseLeaderboardStore) GetUserRankingsFiltered(country string) []server.User {
	scores, err := d.RedisClient.ZRevRangeWithScores(Ctx, leaderboardKey+":"+strings.ToLower(country), 0, -1).Result()
	if err != nil {
		log.Fatal(err)
	}
	var users []server.User
	for _, score := range scores {
		str, _ := score.Member.(string)
		var user server.User
		findOptions := options.FindOne()
		err = d.getUsers().FindOne(Ctx, bson.M{"display_name": str}, findOptions).Decode(&user)
		rank, err := d.RedisClient.ZRevRank(Ctx,  leaderboardKey+":"+strings.ToLower(country), user.DisplayName).Result()
		if err != nil {
			return nil
		}
		user.Rank = int64(rank + 1)
		users = append(users, user)
	}
	return users
}

func (d *DatabaseLeaderboardStore) CreateUserProfile(user server.User) (server.User, error) {
	user.UserId = uuid.New().String()
	d.userLock.Lock()
	defer d.userLock.Unlock()
	var u server.User
	filter := bson.D{{"display_name", user.DisplayName}}
	err := d.getUsers().FindOne(nil, filter).Decode(&u)
	if err == nil {
		return server.User{}, server.UserExistsError
	}
	user.Points = 0
	_, dbError := d.getUsers().InsertOne(Ctx, user)
	if dbError != nil {
		return server.User{}, dbError
	}
	member := &redis.Z{
		Score: float64(user.Points),
		Member: user.DisplayName,
	}
	d.RedisClient.ZAdd(Ctx, leaderboardKey, member)
	d.RedisClient.ZAdd(Ctx, leaderboardKey+":"+strings.ToLower(user.Country), member)
	rank, err := d.RedisClient.ZRevRank(Ctx, leaderboardKey, user.DisplayName).Result()
	if err != nil {
		return server.User{}, ErrRankAcquire
	}
	user.Rank = rank + 1
	return user, nil
}

func (d *DatabaseLeaderboardStore) GetUserProfile(userId string) (server.User, error) {
	var u server.User
	filter := bson.M{ "_id": bson.M{"$eq": userId}}
	if err := d.getUsers().FindOne(nil, filter).Decode(&u); err != nil {
		return server.User{}, dbError
	}
	rank, err := d.RedisClient.ZRevRank(Ctx, leaderboardKey, u.UserId).Result()
	if err != nil {
		return server.User{}, ErrRankAcquire
	}
	u.Rank = rank + 1
	return u, nil
}

func (d *DatabaseLeaderboardStore) SubmitUserScore(score server.Score) (server.Score, error) {
	var u server.User
	score.TimeStamp = time.Now().String()
	update := bson.D{
		{"$inc", bson.D{
			{"points", score.Score}}},
			{"$set", bson.D{{"last_score_timestamp", score.TimeStamp}}}}

	err := d.getUsers().FindOneAndUpdate(nil, bson.M{"_id" : score.UserId}, update).Decode(&u)
	if err != nil {
		return score, server.NoUserPresentError
	}
	_, err = d.RedisClient.ZIncrBy(Ctx,leaderboardKey, float64(score.Score), u.DisplayName).Result()
	if err != nil {
		return score, server.NoUserPresentError
	}
	_, err = d.RedisClient.ZIncrBy(Ctx,leaderboardKey+":"+strings.ToLower(u.Country), float64(score.Score), u.DisplayName).Result()
	if err != nil {
		return score, server.NoUserPresentError
	}

	return score, nil
}

func (d *DatabaseLeaderboardStore) CreateUserProfiles(submission server.Submission) error {
	userSize := submission.SubmissionSize
	d.userLock.Lock()
	defer d.userLock.Unlock()
	count, err := d.getUsers().CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	var userList []interface{}
	for index := 1; index <= userSize ; index++ {
		userList = append(userList, server.User{
			Rank: int64(count) + int64(index),
			UserId:      uuid.New().String(),
			DisplayName: randstr.String(10),
			Country:     getRandomEntry(countryList)})
	}
	d.getUsers().InsertMany(context.TODO(), userList)
	return nil
}
func (d *DatabaseLeaderboardStore) CreateScoreSubmissions(submission server.Submission) error {

	return nil
}

func (d *DatabaseLeaderboardStore) getUsers() *mongo.Collection {
	return d.MongoClient.Database("leaderboard").Collection("users")
}

func NewDatabaseLeaderboardStore(config server.ConfigurationType) *DatabaseLeaderboardStore {
	return &DatabaseLeaderboardStore{
		nil,
		nil,
		config.MongoUri,
		config.RedisUri,
		&sync.Mutex{},
		&sync.Mutex{}}
}

func (d *DatabaseLeaderboardStore) InitializeConnection() func() {
	clientOptions := options.Client().ApplyURI(d.mongoUri)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	d.MongoClient = client
	mod := mongo.IndexModel{
		Keys: bson.M{
			"display_name": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}
	d.getUsers().Indexes().CreateOne(Ctx, mod)
	redisClient := redis.NewClient(&redis.Options{
		Addr: d.redisUri,
		Password: "",
		DB: 0,
	})
	if err := redisClient.Ping(Ctx).Err(); err != nil {
		panic(err)
	}
	d.RedisClient = redisClient
	closeConnection := func() {
		_ = d.MongoClient.Disconnect(ctx)
	}
	return closeConnection
}

func (d *DatabaseLeaderboardStore) InitializeRedisCache() {
	pipeline := []bson.M{
		{
			"$sort": bson.M{
				"points": -1,
			},
		},
	}

	d.RedisClient.FlushAll(Ctx)

	cursor, err := d.getUsers().Aggregate(Ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	var users []server.User
	cursor.All(Ctx, &users)

	for _, user := range users {
		_, err := d.RedisClient.ZAdd(Ctx, leaderboardKey, &redis.Z{
			Score:  float64(user.Points),
			Member: user.DisplayName,
		}).Result()
		if err != nil {
			log.Fatal(err)
		}
		_, err = d.RedisClient.ZAdd(Ctx, leaderboardKey+":"+strings.ToLower(user.Country), &redis.Z{
			Score:  float64(user.Points),
			Member: user.DisplayName,
		}).Result()

		if err != nil {
			log.Fatal(err)
		}
	}
}
