package mcptools

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tmdbAPIKeyFromEnv(t *testing.T) string {
	t.Helper()
	key := os.Getenv("TMDB_API_KEY")
	if key == "" {
		t.Skip("TMDB_API_KEY environment variable not set")
	}
	return key
}

func TestSearchMovies(t *testing.T) {
	key := tmdbAPIKeyFromEnv(t)
	tmdb := NewTMDB(key, "en-US")

	tests := []struct {
		name  string
		input TMDBSearchMovieInput
	}{
		{
			name: "Search for The Matrix",
			input: TMDBSearchMovieInput{
				Name: "The Matrix",
			},
		},
		{
			name: "Search for The Matrix with year",
			input: TMDBSearchMovieInput{
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
	key := tmdbAPIKeyFromEnv(t)
	tmdb := NewTMDB(key, "en-US")
	// Year is wrong
	result, err := tmdb.searchMovies(TMDBSearchMovieInput{Name: "The Matrix", Year: 1990})
	require.NoError(t, err)
	require.Empty(t, result.Results)
}

func TestSearchTVShows(t *testing.T) {
	key := tmdbAPIKeyFromEnv(t)
	tmdb := NewTMDB(key, "en-US")
	result, err := tmdb.searchTVShows(TMDBSearchTVShowInput{Name: "Breaking Bad"})
	require.NoError(t, err)
	require.NotEmpty(t, result.Results)
	assert.Equal(t, "Breaking Bad", result.Results[0].Name)
	assert.Len(t, result.Results[0].Seasons, 6)
}
