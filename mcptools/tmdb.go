package mcptools

import (
	"context"
	"log"
	"strconv"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const TMDB_LIMIT_ACTORS_COUNT = 10

type TMDB struct {
	apiKey   string
	language string
}

func NewTMDBClient(apiKey, language string) *TMDB {
	return &TMDB{
		apiKey:   apiKey,
		language: language,
	}
}

func (s *TMDB) AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_movies",
		Description: "Search for movies on TMDB by given name (required) and year (optional).",
	}, s.SearchMovies)

}

type SearchMovieInput struct {
	Name string `json:"name" jsonschema:"the name of the movie to search for"`
	Year int    `json:"year,omitempty" jsonschema:"(optional) the year of the movie released"`
}

type TMDBActor struct {
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
}

type TMDBMovieItem struct {
	Title            string      `json:"title"`
	OriginalTitle    string      `json:"original_title"`
	OriginalLanguage string      `json:"original_language"`
	Overview         string      `json:"overview"`
	ReleaseDate      string      `json:"release_date"`
	Actors           []TMDBActor `json:"actors"`
}

type SearchMovieOutput struct {
	Results []TMDBMovieItem `json:"results"`
}

func (s *TMDB) searchMovies(input SearchMovieInput) (SearchMovieOutput, error) {
	c, err := tmdb.Init(s.apiKey)
	if err != nil {
		log.Printf("Error initializing TMDB client: %v", err)
		return SearchMovieOutput{}, err
	}

	options := map[string]string{"language": s.language, "include_adult": "true"}
	if input.Year != 0 {
		options["year"] = strconv.Itoa(input.Year)
	}

	searchRes, err := c.GetSearchMovies(input.Name, options)
	if err != nil {
		log.Printf("Error searching movies: %v", err)
		return SearchMovieOutput{}, err
	}

	var results []TMDBMovieItem
	for _, movie := range searchRes.Results {
		movieItem := TMDBMovieItem{
			Title:            movie.Title,
			OriginalTitle:    movie.OriginalTitle,
			OriginalLanguage: movie.OriginalLanguage,
			Overview:         movie.Overview,
			ReleaseDate:      movie.ReleaseDate,
		}

		var actors []TMDBActor
		options := map[string]string{"language": s.language}
		credits, err := c.GetMovieCredits(int(movie.ID), options)
		if err != nil {
			log.Printf("Error getting movie credits: %v", err)
		} else {
			for _, cast := range credits.Cast {
				if len(actors) >= TMDB_LIMIT_ACTORS_COUNT {
					break
				}
				actors = append(actors, TMDBActor{
					Name:         cast.Name,
					OriginalName: cast.OriginalName,
				})
			}
			movieItem.Actors = actors
		}

		results = append(results, movieItem)
	}

	return SearchMovieOutput{Results: results}, nil
}

func (s *TMDB) SearchMovies(
	ctx context.Context, req *mcp.CallToolRequest, input SearchMovieInput) (
	*mcp.CallToolResult, SearchMovieOutput, error) {
	result, err := s.searchMovies(input)
	return nil, result, err
}
