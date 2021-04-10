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

func main() {
	configuration, err := ParseArguments(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	store, closeConnection := datastore.NewDatabaseLeaderboardStore(configuration)
	defer closeConnection()
	var leaderboardServer *server.LeaderboardServer
	leaderboardServer = ConfigureServer(configuration, store)
	log.Fatal(http.ListenAndServe(":5000", leaderboardServer))
}

func ParseArguments(args []string) (server.ConfigurationType, error) {
	if len(args) == 2 && args[1] != serverTypeDev {
		return server.ConfigurationType{}, errorWrongConfiguration
	}
	if len(args) == 3 && args[1] != serverTypeProd {
		return server.ConfigurationType{}, errorWrongConfiguration
	}
	if args[1] == serverTypeProd {
		return server.ConfigurationType{
			Server:     serverTypeProd,
			Connection: args[2],
			Message:    "Initializing Production Server",
		}, nil
	}
	return server.ConfigurationType{
		Server:     serverTypeDev,
		Connection: localUri,
		Message:    "Initializing Development Server",
	}, nil
}

func ConfigureServer(s server.ConfigurationType, store server.LeaderboardStore) *server.LeaderboardServer {
	var leaderboardServer *server.LeaderboardServer
	if s.Server == serverTypeDev {
		leaderboardServer = server.NewLeaderboardServer(store, true)
	} else {
		leaderboardServer = server.NewLeaderboardServer(store, false)
	}
	fmt.Println(s.Message)
	return leaderboardServer
}
