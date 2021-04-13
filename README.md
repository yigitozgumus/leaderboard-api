# leaderboard

This project is a demo application for GJG case. The application is an api that tracks a scoreboard. The in memory implementation
is done to test the server endpoint with a more basic approach and get use to the composition styles of Go. The Users are kept 
in a MongoDB instance which can be connected either via an .env file or with an appropriate variable. 

The leaderboard itself and score tracking is done with Redis using an ordered set. Connection principle is same with mongodb.
At the initialization phase, server creates the connections to both mongodb and redis then dumps the users' score to redis to initialize
the leaderboard. Example env file is added for env variable reference.

The application can be accessed from [here](https://shielded-basin-78769.herokuapp.com/leaderboard).

## Project Structure 

```bash
.
├── Dockerfile
├── LICENSE
├── README.md
├── datastore
│ ├── database_leaderboard_store.go
│ ├── gotemplate_RankingMap.go
│ ├── in_memory_leaderboard_store.go
│ ├── in_memory_ranking_leaderboard_store.go
│ └── util.go
├── go.mod
├── go.sum
├── main.go
└── server
    ├── server.go
    └── server_test.go

```


## Endpoints 

This section briefly explains what each endpoint retrieves. There are definitely certain improvements to be made for the 
general implementation in terms of the performance.

### GET /leaderboard

Returns the current leaderboard. If the userbase grows really big, then this approach becomes very inefficient due
to the amount of data it needs to deliver. Possible improvements include adding a limit factor, pagination or 
restricting leaderboard to a user specific range. 

```json

[
  {
    "user_id": "efd098ef-57f8-4186-a33a-9789ff1ffe84",
    "display_name": "Cansu",
    "points": 500,
    "rank": 1,
    "country": "fr",
    "LastScoreTimeStamp": ""
  },
  {
    "user_id": "774bac2d-c4de-487b-b516-73709c8a4e8b",
    "display_name": "Yigit",
    "points": 450,
    "rank": 2,
    "country": "de",
    "LastScoreTimeStamp": ""
  },
  {
    "user_id": "cd1e3ed9-628b-4057-a242-c7ef9f224020",
    "display_name": "Aysu",
    "points": 250,
    "rank": 3,
    "country": "fr",
    "LastScoreTimeStamp": ""
  },
  {
    "user_id": "bfdf2b7c-f98a-4756-a92b-807dd4135fc5",
    "display_name": "Mert",
    "points": 120,
    "rank": 4,
    "country": "au",
    "LastScoreTimeStamp": ""
  }
]
```

### GET /leaderboard/{country}

Just like the above endpoint but with a country filter. Aforementioned improvements can also be applied here.

```json

[
  {
    "user_id": "22f6b394-ec26-43d7-a92d-e6a321f62be4",
    "display_name": "Yigit",
    "points": 500,
    "rank": 1,
    "country": "de",
    "LastScoreTimeStamp": "2021-04-12 19:24:25.653688 +0300 +03 m=+24.799103861"
  },
  {
    "user_id": "194c2148-49a4-459c-b600-12eef095a014",
    "display_name": "zBRJucIaep",
    "points": 120,
    "rank": 2,
    "country": "de",
    "LastScoreTimeStamp": ""
  }
]

```

### POST /leaderboard/score/submit

This endpoint is used for score submission. Embedding user id to a score is not a real use case but added for
to simplify the implementation

```json
{
	"score": 500,
	"user_id": "22f6b394-ec26-43d7-a92d-e6a321f62be4"
}
```

```json
{
  "submission": {
    "score": 500,
    "user_id": "22f6b394-ec26-43d7-a92d-e6a321f62be4",
    "time_stamp": "2021-04-12 19:24:25.653688 +0300 +03 m=+24.799103861"
  }
}
```

### GET /leaderboard/user/profile/{id}

Returns the user profile with her/his current rank and total points.

```json
{
  "user_id": "c0209197-c900-4e2e-9dbe-f4f1233a18bf",
  "display_name": "Yigit",
  "points": 650,
  "rank": 1,
  "country": "tr"
}
```

### POST /leaderboard/user/create

 Used for creating new users, returns the added user information.

```json
{
	"display_name": "Yigit",
	"country": "de"
}
```

```json
{
  "user_id": "3d2bfa42-6a39-47f6-99eb-e4d2e7e5f1b1",
  "display_name": "Yigit",
  "points": 0,
  "rank": 2,
  "country": "de",
  "LastScoreTimeStamp": ""
}
```

### POST /leaderboard/test/create/users (dev only)

This endpoint should only used for testing. There are libraries for this type of mock data injection but I wanted to 
fiddle around more with the api itself. Not the best approach but it'll do.

```json
{
	"submission_size": 5000
}
```

### POST /leaderboard/test/submit/scores (dev only)

This endpoint gets all the current users from the current state of the leaderboard and randomly assings scores. If the 
submission size is too big, it may return 503 since there it literally waits all the score submissions before finishing.

```json
{
	"submission_size": 10000,
	"max_score": 250,
	"min_score": 10
}
```


