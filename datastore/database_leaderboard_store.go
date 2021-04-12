package datastore

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
	"github.com/yigitozgumus/leaderboard-api/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"math"
	"sync"
	"time"
)

var dbError = errors.New("error with database")
var ErrScoreSubmission = errors.New("Score submission failure")

type DatabaseLeaderboardStore struct {
	client         *mongo.Client
	ctx            context.Context
	connection     string
	userLock       *sync.Mutex
	scoreLock      *sync.Mutex
}

type Ranking struct {
	Score float64 `bson:"score"`
	CurrentUsers []UserRank `bson:"current_users"`
}

type UserRank struct {
	User string `bson:"user_id"`
}

func (d *DatabaseLeaderboardStore) GetUserRankings() []server.User {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := d.client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (d *DatabaseLeaderboardStore) GetUserRankingsFiltered(country string) []server.User {
	// FIXME
	return nil
}

func (d *DatabaseLeaderboardStore) CreateUserProfile(user server.User) error {
	user.UserId = uuid.New().String()
	d.userLock.Lock()
	defer d.userLock.Unlock()
	filter := bson.D{{"display_name", user.DisplayName}}
	_ = d.getUsers().FindOne(nil, filter)
	count, err := d.getUsers().CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	user.Rank = int64(count) + 1
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, dbError := d.getUsers().InsertOne(ctx, user)
	if dbError != nil {
		return dbError
	}
	find := d.getRankings().FindOne(nil, bson.M{"score": 0})
	if find.Err() != nil {
		d.getRankings().InsertOne(nil,Ranking{0, []UserRank{UserRank{user.UserId}}})
	} else {
		filter := bson.M{
			"score": bson.M{
				"$eq": 0,
			},
		}
		update := bson.M{
			"$push": bson.M{
				"current_users": bson.M{"user_id": user.UserId}}}
		_ = d.getRankings().FindOneAndUpdate(nil, filter, update)
	}

	return nil
}

func (d *DatabaseLeaderboardStore) GetUserProfile(userId string) (server.User, error) {
	var u server.User
	filter := bson.M{ "_id": bson.M{"$eq": userId}}
	err := d.getUsers().FindOne(nil, filter).Decode(&u)
	if err != nil {
		return server.User{}, dbError
	}
	return u, nil
	// FIXME update user's current ranking
}

func (d *DatabaseLeaderboardStore) SubmitUserScore(score server.Score) (server.Score, error) {
	score.TimeStamp = time.Now().String()
	// check if the user is present
	find := d.getUsers().FindOne(nil, bson.M{"_id" : score.UserId})
	if find.Err() != nil {
		return score, server.NoUserPresentError
	}
	d.scoreLock.Lock()
	defer d.scoreLock.Unlock()
	var u server.User
	find.Decode(&u)
	currentScore := u.Points
	newScore := currentScore + score.Score
	newScore = math.Round(newScore*100) / 100
	// check if current score of the player is present in the leaderboard
	find = d.getRankings().FindOne(nil, bson.M{"score": currentScore})
	if find.Err() != nil {
		return score, errors.New("current score should have a entry")
	}
	err := d.removeOldScoreOfUser(currentScore, score.UserId)
	if err != nil {
		return score, ErrScoreSubmission
	}
	err = d.updateNewScore(newScore, score.UserId)
	if err != nil {
		d.getRankings().InsertOne(nil,Ranking{newScore, []UserRank{UserRank{score.UserId}}})
	}
	err = d.updateUserTotalPoints(newScore, score.UserId)
	if err != nil {
		return score, ErrScoreSubmission
	}
	return score, nil
}

func (d *DatabaseLeaderboardStore) removeOldScoreOfUser(currentScore float64, userId string ) error {
	filter := bson.M{ "score": bson.M{ "$eq": currentScore }}
	update := bson.M{
		"$pull": bson.M{
			"current_users" : bson.M {
				"user_id": bson.M{"$in": bson.A{userId}}}},
	}
	res := d.getRankings().FindOneAndUpdate(nil, filter, update)
	return res.Err()
}

func (d *DatabaseLeaderboardStore) updateNewScore(newScore float64, userId string) error {
	filter := bson.M{ "score": bson.M{ "$eq": newScore }}
	update := bson.M{ "$push": bson.M{ "current_users.user_id": userId}}
	res := d.getRankings().FindOneAndUpdate(nil, filter, update)
	return res.Err()
}

func (d *DatabaseLeaderboardStore) updateUserTotalPoints(updatedScore float64, userId string) error {
	filter := bson.M { "_id": bson.M { "$eq": userId } }
	update := bson.M { "$set" : bson.M { "points": updatedScore }}
	res := d.getUsers().FindOneAndUpdate(nil , filter, update)
	return res.Err()
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
	d.client.Database("leaderboard").Collection("users").InsertMany(context.TODO(), userList)
	return nil
}
func (d *DatabaseLeaderboardStore) CreateScoreSubmissions(submission server.Submission) error {
	// FIXME
	return nil
}

func (d *DatabaseLeaderboardStore) getUsers() *mongo.Collection {
	return d.client.Database("leaderboard").Collection("users")
}

func (d *DatabaseLeaderboardStore) getRankings() *mongo.Collection {
	return d.client.Database("leaderboard").Collection("rankings")
}

func NewDatabaseLeaderboardStore(config server.ConfigurationType) *DatabaseLeaderboardStore {
	return &DatabaseLeaderboardStore{
		nil,
		nil,
		config.Connection,
		&sync.Mutex{},
		&sync.Mutex{}}
}

func (d *DatabaseLeaderboardStore) InitializeConnection() func() {
	clientOptions := options.Client().ApplyURI(d.connection)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	d.client = client
	closeConnection := func() {
		_ = d.client.Disconnect(ctx)
	}
	return closeConnection
}
