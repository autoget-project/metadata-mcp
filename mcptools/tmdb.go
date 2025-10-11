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

func NewTMDB(apiKey, language string) *TMDB {
	return &TMDB{
		apiKey:   apiKey,
		language: language,
	}
}

func (s *TMDB) AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_movies",
		Description: "Search for movies on TMDB by given name (required) and year (optional).",
	}, s.searchMoviesTool)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_tv_shows",
		Description: "Search for tv shows on TMDB by given name.",
	}, s.searchTVShowsTool)
}

type TMDBSearchMovieInput struct {
	Name string `json:"name" jsonschema:"the name of the movie or tv show to search for"`
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

func (s *TMDB) searchMovies(input TMDBSearchMovieInput) (SearchMovieOutput, error) {
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
				if cast.KnownForDepartment != "Acting" {
					continue
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

func (s *TMDB) searchMoviesTool(
	ctx context.Context, req *mcp.CallToolRequest, input TMDBSearchMovieInput) (
	*mcp.CallToolResult, SearchMovieOutput, error) {
	result, err := s.searchMovies(input)
	return nil, result, err
}

type TMDBSearchTVShowInput struct {
	Name string `json:"name" jsonschema:"the name of the movie or tv show to search for"`
}

type TMDBTVShowSeason struct {
	Name         string `json:"name"`
	SeasonNumber int    `json:"season_number"`
	EpisodeCount int    `json:"episode_count"`
	AirDate      string `json:"air_date"`
}

type TMDBTVShowItem struct {
	Name             string             `json:"name"`
	OriginalName     string             `json:"original_name"`
	OriginalLanguage string             `json:"original_language"`
	Overview         string             `json:"overview"`
	FirstAirDate     string             `json:"first_air_date"`
	Actors           []TMDBActor        `json:"actors"`
	Seasons          []TMDBTVShowSeason `json:"seasons"`
}

type SearchTVShowOutput struct {
	Results []TMDBTVShowItem `json:"results"`
}

func (s *TMDB) searchTVShows(input TMDBSearchTVShowInput) (SearchTVShowOutput, error) {
	c, err := tmdb.Init(s.apiKey)
	if err != nil {
		log.Printf("Error initializing TMDB client: %v", err)
		return SearchTVShowOutput{}, err
	}

	options := map[string]string{"language": s.language, "include_adult": "true"}
	searchRes, err := c.GetSearchTVShow(input.Name, options)
	if err != nil {
		log.Printf("Error searching tv shows: %v", err)
		return SearchTVShowOutput{}, err
	}

	var results []TMDBTVShowItem
	for _, tvShow := range searchRes.Results {
		item := TMDBTVShowItem{
			Name:             tvShow.Name,
			OriginalName:     tvShow.OriginalName,
			OriginalLanguage: tvShow.OriginalLanguage,
			Overview:         tvShow.Overview,
			FirstAirDate:     tvShow.FirstAirDate,
		}

		// get seasons
		options := map[string]string{"language": s.language, "append_to_response": "credits"}
		details, err := c.GetTVDetails(int(tvShow.ID), options)
		if err != nil {
			log.Printf("Error getting tv details: %v", err)
		} else {
			for _, season := range details.Seasons {
				item.Seasons = append(item.Seasons, TMDBTVShowSeason{
					Name:         season.Name,
					SeasonNumber: season.SeasonNumber,
					EpisodeCount: season.EpisodeCount,
					AirDate:      season.AirDate,
				})
			}

			for _, cast := range details.Credits.Cast {
				if len(item.Actors) >= TMDB_LIMIT_ACTORS_COUNT {
					break
				}
				if cast.KnownForDepartment != "Acting" {
					continue
				}
				item.Actors = append(item.Actors, TMDBActor{
					Name:         cast.Name,
					OriginalName: cast.OriginalName,
				})
			}
		}

		results = append(results, item)
	}

	return SearchTVShowOutput{Results: results}, nil
}

func (s *TMDB) searchTVShowsTool(
	ctx context.Context, req *mcp.CallToolRequest, input TMDBSearchTVShowInput) (
	*mcp.CallToolResult, SearchTVShowOutput, error) {
	result, err := s.searchTVShows(input)
	return nil, result, err
}
