# leaderboard

This project is a demo application for GJG case. It will be a leaderboard api 
that will return users ranked by their points.

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


## How to build

for command line application (in progress)

```bash
go build ./cmd/cli/
go run ./cmd/cli/main.go
```

For webserver

```bash
go build ./cmd/webserver
go run ./cmd/webserver/main.go
```