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

func TestFindByIMDB(t *testing.T) {
	key := tmdbAPIKeyFromEnv(t)
	tmdb := NewTMDB(key, "en-US")

	tests := []struct {
		name         string
		imdbID       string
		expectType   string
		expectName   string
		checkDetails func(t *testing.T, result TMDBFindByIMDBOutput)
	}{
		{
			name:       "Movie: 12 Angry Men",
			imdbID:     "tt0050083",
			expectType: "movie",
			expectName: "12 Angry Men",
			checkDetails: func(t *testing.T, result TMDBFindByIMDBOutput) {
				require.NotEmpty(t, result.MovieResults)
				assert.Equal(t, "12 Angry Men", result.MovieResults[0].Title)
				assert.NotEmpty(t, result.MovieResults[0].Overview)
				assert.NotEmpty(t, result.MovieResults[0].Actors)
			},
		},
		{
			name:       "TV Show: Only Murders in the Building",
			imdbID:     "tt11691774",
			expectType: "tv",
			expectName: "Only Murders in the Building",
			checkDetails: func(t *testing.T, result TMDBFindByIMDBOutput) {
				require.NotEmpty(t, result.TVResults)
				assert.Equal(t, "Only Murders in the Building", result.TVResults[0].Name)
				assert.NotEmpty(t, result.TVResults[0].Seasons)
				assert.NotEmpty(t, result.TVResults[0].Actors)
			},
		},
		{
			name:       "TV Episode: Only Murders in the Building episode",
			imdbID:     "tt33333196",
			expectType: "tv",
			expectName: "Only Murders in the Building",
			checkDetails: func(t *testing.T, result TMDBFindByIMDBOutput) {
				require.NotEmpty(t, result.TVResults)
				assert.Equal(t, "Only Murders in the Building", result.TVResults[0].Name)
				assert.NotEmpty(t, result.TVResults[0].Seasons)
			},
		},
		{
			name:       "Person: Steve Martin",
			imdbID:     "nm0000188",
			expectType: "person",
			expectName: "Steve Martin",
			checkDetails: func(t *testing.T, result TMDBFindByIMDBOutput) {
				require.NotEmpty(t, result.PersonResults)
				assert.Equal(t, "Steve Martin", result.PersonResults[0].Name)
				assert.Greater(t, result.PersonResults[0].Popularity, float32(0))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tmdb.findByIMDB(TMDBFindByIMDBInput{IMDBID: tt.imdbID})
			require.NoError(t, err)

			switch tt.expectType {
			case "movie":
				require.NotEmpty(t, result.MovieResults)
				assert.Equal(t, tt.expectName, result.MovieResults[0].Title)
			case "tv":
				require.NotEmpty(t, result.TVResults)
				assert.Equal(t, tt.expectName, result.TVResults[0].Name)
			case "person":
				require.NotEmpty(t, result.PersonResults)
				assert.Equal(t, tt.expectName, result.PersonResults[0].Name)
			}

			if tt.checkDetails != nil {
				tt.checkDetails(t, result)
			}
		})
	}
}
