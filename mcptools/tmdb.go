package mcptools

import (
	"context"
	"log"
	"strconv"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const tmdbLimitActorsCount = 10

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
		Description: "Searches for movies on The Movie Database (TMDB) by name (required) and optional release year.",
	}, s.searchMoviesTool)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_tv_shows",
		Description: "Searches for TV shows on The Movie Database (TMDB) by name.",
	}, s.searchTVShowsTool)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "find_by_imdb_id",
		Description: "Finds content on TMDB by IMDB ID using external source lookup.",
	}, s.findByIMDBTool)
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
	Actors           []TMDBActor `json:"actors,omitempty"`
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
				if len(actors) >= tmdbLimitActorsCount {
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
	Actors           []TMDBActor        `json:"actors,omitempty"`
	Seasons          []TMDBTVShowSeason `json:"seasons,omitempty"`
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
				if len(item.Actors) >= tmdbLimitActorsCount {
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

type TMDBFindByIMDBInput struct {
	IMDBID string `json:"imdb_id" jsonschema:"the IMDB ID to search for (e.g., 'tt0111161')"`
}

type TMDBPerson struct {
	Name         string  `json:"name"`
	OriginalName string  `json:"original_name"`
	Birthday     string  `json:"birthday,omitempty"`
	PlaceOfBirth string  `json:"place_of_birth,omitempty"`
	Deathday     string  `json:"deathday,omitempty"`
	Popularity   float32 `json:"popularity,omitempty"`
}

type TMDBFindByIMDBOutput struct {
	MovieResults    []TMDBMovieItem `json:"movie_results,omitempty"`
	TVResults       []TMDBTVShowItem `json:"tv_results,omitempty"`
	PersonResults   []TMDBPerson     `json:"person_results,omitempty"`
}

func (s *TMDB) findByIMDB(input TMDBFindByIMDBInput) (TMDBFindByIMDBOutput, error) {
	c, err := tmdb.Init(s.apiKey)
	if err != nil {
		log.Printf("Error initializing TMDB client: %v", err)
		return TMDBFindByIMDBOutput{}, err
	}

	options := map[string]string{
		"language":        s.language,
		"external_source": "imdb_id",
	}

	findResult, err := c.GetFindByID(input.IMDBID, options)
	if err != nil {
		log.Printf("Error finding by IMDB ID: %v", err)
		return TMDBFindByIMDBOutput{}, err
	}

	result := TMDBFindByIMDBOutput{}

	// Handle movie results
	for _, movie := range findResult.MovieResults {
		detailOptions := map[string]string{"language": s.language, "append_to_response": "credits"}
		details, err := c.GetMovieDetails(int(movie.ID), detailOptions)
		if err != nil {
			log.Printf("Error getting movie details: %v", err)
			continue
		}

		movieItem := TMDBMovieItem{
			Title:            details.Title,
			OriginalTitle:    details.OriginalTitle,
			OriginalLanguage: details.OriginalLanguage,
			Overview:         details.Overview,
			ReleaseDate:      details.ReleaseDate,
		}

		var actors []TMDBActor
		for _, cast := range details.Credits.Cast {
			if len(actors) >= tmdbLimitActorsCount {
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
		result.MovieResults = append(result.MovieResults, movieItem)
	}

	// Handle TV results
	for _, tvShow := range findResult.TvResults {
		detailOptions := map[string]string{"language": s.language, "append_to_response": "credits"}
		details, err := c.GetTVDetails(int(tvShow.ID), detailOptions)
		if err != nil {
			log.Printf("Error getting tv details: %v", err)
			continue
		}

		tvItem := TMDBTVShowItem{
			Name:             details.Name,
			OriginalName:     details.OriginalName,
			OriginalLanguage: details.OriginalLanguage,
			Overview:         details.Overview,
			FirstAirDate:     details.FirstAirDate,
		}

		// get seasons
		for _, season := range details.Seasons {
			tvItem.Seasons = append(tvItem.Seasons, TMDBTVShowSeason{
				Name:         season.Name,
				SeasonNumber: season.SeasonNumber,
				EpisodeCount: season.EpisodeCount,
				AirDate:      season.AirDate,
			})
		}

		for _, cast := range details.Credits.Cast {
			if len(tvItem.Actors) >= tmdbLimitActorsCount {
				break
			}
			if cast.KnownForDepartment != "Acting" {
				continue
			}
			tvItem.Actors = append(tvItem.Actors, TMDBActor{
				Name:         cast.Name,
				OriginalName: cast.OriginalName,
			})
		}
		result.TVResults = append(result.TVResults, tvItem)
	}

	// Handle TV episode results (find the TV series)
	for _, episode := range findResult.TvEpisodeResults {
		detailOptions := map[string]string{"language": s.language, "append_to_response": "credits"}
		details, err := c.GetTVDetails(int(episode.ShowID), detailOptions)
		if err != nil {
			log.Printf("Error getting tv series details for episode: %v", err)
			continue
		}

		tvItem := TMDBTVShowItem{
			Name:             details.Name,
			OriginalName:     details.OriginalName,
			OriginalLanguage: details.OriginalLanguage,
			Overview:         details.Overview,
			FirstAirDate:     details.FirstAirDate,
		}

		// get seasons
		for _, season := range details.Seasons {
			tvItem.Seasons = append(tvItem.Seasons, TMDBTVShowSeason{
				Name:         season.Name,
				SeasonNumber: season.SeasonNumber,
				EpisodeCount: season.EpisodeCount,
				AirDate:      season.AirDate,
			})
		}

		for _, cast := range details.Credits.Cast {
			if len(tvItem.Actors) >= tmdbLimitActorsCount {
				break
			}
			if cast.KnownForDepartment != "Acting" {
				continue
			}
			tvItem.Actors = append(tvItem.Actors, TMDBActor{
				Name:         cast.Name,
				OriginalName: cast.OriginalName,
			})
		}
		result.TVResults = append(result.TVResults, tvItem)
	}

	// Handle TV season results (find the TV series)
	for _, season := range findResult.TvSeasonResults {
		detailOptions := map[string]string{"language": s.language, "append_to_response": "credits"}
		details, err := c.GetTVDetails(int(season.ShowID), detailOptions)
		if err != nil {
			log.Printf("Error getting tv series details for season: %v", err)
			continue
		}

		tvItem := TMDBTVShowItem{
			Name:             details.Name,
			OriginalName:     details.OriginalName,
			OriginalLanguage: details.OriginalLanguage,
			Overview:         details.Overview,
			FirstAirDate:     details.FirstAirDate,
		}

		// get seasons
		for _, s := range details.Seasons {
			tvItem.Seasons = append(tvItem.Seasons, TMDBTVShowSeason{
				Name:         s.Name,
				SeasonNumber: s.SeasonNumber,
				EpisodeCount: s.EpisodeCount,
				AirDate:      s.AirDate,
			})
		}

		for _, cast := range details.Credits.Cast {
			if len(tvItem.Actors) >= tmdbLimitActorsCount {
				break
			}
			if cast.KnownForDepartment != "Acting" {
				continue
			}
			tvItem.Actors = append(tvItem.Actors, TMDBActor{
				Name:         cast.Name,
				OriginalName: cast.OriginalName,
			})
		}
		result.TVResults = append(result.TVResults, tvItem)
	}

	// Handle person results
	for _, person := range findResult.PersonResults {
		personItem := TMDBPerson{
			Name:         person.Name,
			OriginalName: person.Name, // Person API doesn't have OriginalName
			Popularity:   person.Popularity,
		}
		result.PersonResults = append(result.PersonResults, personItem)
	}

	return result, nil
}

func (s *TMDB) findByIMDBTool(
	ctx context.Context, req *mcp.CallToolRequest, input TMDBFindByIMDBInput) (
	*mcp.CallToolResult, TMDBFindByIMDBOutput, error) {
	result, err := s.findByIMDB(input)
	return nil, result, err
}
