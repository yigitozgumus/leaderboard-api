# leaderboard

This project is a demo application for GJG case. It will be a leaderboard api 
that will return users ranked by their points.

## Project Structure 

```bash
.
├── LICENSE
├── README.md
├── cmd
│ ├── cli
│ └── webserver
│     └── main.go
├── db
│ └── in_memory_leaderboard_store.go
├── go.mod
├── pbcopy
├── server
│ ├── server.go
│ └── server_test.go
└── webserver
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