package main

import (
	"errors"
	"fmt"
	"github.com/yigitozgumus/leaderboard-api/datastore"
	"github.com/yigitozgumus/leaderboard-api/server"
	"log"
	"net/http"
	"os"
)

const localUri = "mongodb://localhost:27017"
const serverTypeDev = "dev"
const serverTypeProd = "prod"

var errorWrongConfiguration = errors.New("wrong server configuration")

type ConfigurationType struct {
	server     string
	connection string
	message    string
}

func main() {
	serverType, err := ParseArguments(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	store := datastore.NewInMemoryRankingStore()
	var leaderboardServer *server.LeaderboardServer
	leaderboardServer = ConfigureServer(serverType, store)
	log.Fatal(http.ListenAndServe(":5000", leaderboardServer))
}

func ParseArguments(args []string) (ConfigurationType, error) {
	if len(args) == 2 && args[1] != serverTypeDev {
		return ConfigurationType{}, errorWrongConfiguration
	}
	if len(args) == 3 && args[1] != serverTypeProd {
		return ConfigurationType{}, errorWrongConfiguration
	}
	if args[1] == serverTypeProd {
		return ConfigurationType{
			server:     serverTypeProd,
			connection: args[2],
			message: "Initializing Production Server",
		}, nil
	}
	return ConfigurationType{
		server:     serverTypeDev,
		connection: localUri,
		message: "Initializing Development Server",
	}, nil
}

func ConfigureServer(s ConfigurationType, store server.LeaderboardStore) *server.LeaderboardServer {
	var leaderboardServer *server.LeaderboardServer
	if s.server == serverTypeDev {
		leaderboardServer = server.NewLeaderboardServer(store, true)
	} else {
		leaderboardServer = server.NewLeaderboardServer(store, false)
	}
	fmt.Println(s.message)
	return leaderboardServer
}
