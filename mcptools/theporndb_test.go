package mcptools

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func tpdbTokenFromEnv(t *testing.T) string {
	t.Helper()
	token := os.Getenv("TPDB_API_TOKEN")
	if token == "" {
		t.Skip("TPDB_API_TOKEN environment variable not set")
	}
	return token
}

func TestSearchTPDBVideo(t *testing.T) {
	token := tpdbTokenFromEnv(t)
	tpdb := NewThePornDB(token)
	got, err := tpdb.searchTPDBVideos(t.Context(), TPDBSearchVideosInput{Query: "Long Con"})
	require.NoError(t, err)
	require.NotEmpty(t, got.Results)
}
