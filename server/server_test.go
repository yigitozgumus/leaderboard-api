package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETLeaderboard(t *testing.T) {
	t.Run("it returns 200 on /leaderboard", func(t *testing.T) {
		request := newLeaderboardRequest()
		response := httptest.NewRecorder()
		LeaderboardServer(response, request)
		got := response.Code
		want := http.StatusOK
		if got != want {
			t.Errorf("did not get correct status, got %d, want %d",got, want )
		}
	})
}

// helpers
func newLeaderboardRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet,"/leaderboard", nil )
	return req
}
