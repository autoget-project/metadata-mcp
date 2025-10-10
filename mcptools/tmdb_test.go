package mcptools

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func apiKeyFromEnv(t *testing.T) string {
	t.Helper()
	key := os.Getenv("TMDB_API_KEY")
	if key == "" {
		t.Skip("TMDB_API_KEY environment variable not set")
	}
	return key
}

func TestSearchMovies(t *testing.T) {
	key := apiKeyFromEnv(t)
	tmdb := NewTMDBClient(key, "en-US")

	tests := []struct {
		name  string
		input SearchMovieInput
	}{
		{
			name: "Search for The Matrix",
			input: SearchMovieInput{
				Name: "The Matrix",
			},
		},
		{
			name: "Search for The Matrix with year",
			input: SearchMovieInput{
				Name: "The Matrix",
				Year: 2000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tmdb.searchMovies(tt.input)
			require.NoError(t, err)
			require.NotEmpty(t, result.Results)
		})
	}
}

func TestSearchMoviesNotExists(t *testing.T) {
	key := apiKeyFromEnv(t)
	tmdb := NewTMDBClient(key, "en-US")
	result, err := tmdb.searchMovies(SearchMovieInput{Name: "The Matrix", Year: 1990})
	require.NoError(t, err)
	require.Empty(t, result.Results)
}
