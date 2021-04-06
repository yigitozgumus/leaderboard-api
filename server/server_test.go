package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETLeaderboard(t *testing.T) {
	server := &LeaderboardServer{}
	t.Run("server is running", func(t *testing.T) {
		request := newLeaderboardRequest("/leaderboard")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		got := response.Code
		want := http.StatusOK
		assertStatus(t, got, want)
	})
	t.Run("malformed endpoint prefix returns 404", func(t *testing.T) {
		request := newLeaderboardRequest("/leaderboar")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		got := response.Code
		want := http.StatusNotFound
		assertStatus(t, got, want)
	})
}

func TestGetLeaderboardFiltered(t *testing.T) {
	// FIXME
}

func TestPOSTUserCreate(t *testing.T) {
	// FIXME
}

func TestPOSTScoreSubmit(t *testing.T) {
	// FIXME
}

func TestGETUserProfile(t *testing.T) {
	// FIXME
}



// helpers
func newLeaderboardRequest(prefix string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet,prefix, nil )
	return req
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d",got, want )
	}
}
