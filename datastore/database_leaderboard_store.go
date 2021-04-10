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
	"sync"
	"time"
)

var dbError = errors.New("error with database")

type DatabaseLeaderboardStore struct {
	client         *mongo.Client
	ctx            context.Context
	connection     string
	userLock       *sync.Mutex
	scoreLock      *sync.Mutex
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
	user.Rank = uint64(count) + 1
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, dbError := d.getUsers().InsertOne(ctx, user)
	if dbError != nil {
		return dbError
	}
	// TODO add ranking information

	return nil
}
func (d *DatabaseLeaderboardStore) GetUserProfile(userId string) (server.User, error) {
	return server.User{}, nil
}

func (d *DatabaseLeaderboardStore) SubmitUserScore(score server.Score) (server.Score, error) {
	score.TimeStamp = time.Now().String()
	find := d.getUsers().FindOne(nil, bson.M{"_id" : score.UserId})
	if find.Err() != nil {
		return score, server.NoUserPresentError
	}
	d.scoreLock.Lock()
	defer d.scoreLock.Unlock()
	var u server.User
	find.Decode(&u)
	currentScore := u.Points
	

	return server.Score{}, nil
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
			Rank: uint64(count) + uint64(index),
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
