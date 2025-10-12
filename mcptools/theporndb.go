package mcptools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	tpdbLimitVideoPerType = 10
	tpdbSearchMovieURL    = "https://api.theporndb.net/movies"
	tpdbSearchSceneURL    = "https://api.theporndb.net/scenes"
)

type ThePornDB struct {
	apiToken string
}

func NewThePornDB(apiToken string) *ThePornDB {
	return &ThePornDB{
		apiToken: apiToken,
	}
}

func (s *ThePornDB) AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_porn",
		Description: "Search Non-Japanese Porns",
	}, s.searchTPDBVideosTool)
}

type tpdbActorInfo struct {
	ID        string  `json:"id"`
	NumericID uint64  `json:"_id"`
	Slug      string  `json:"slug"`
	Name      string  `json:"name"`
	Bio       string  `json:"bio"`
	Rating    float32 `json:"rating"`
	IsParent  bool    `json:"is_parent"`
	Extras    struct {
		Gender            string `json:"gender"`
		Birthday          string `json:"birthday"`
		BirthdayTimestamp int    `json:"birthday_timestamp"`
		Birthplace        string `json:"birthplace"`
		BirthplaceCode    string `json:"birthplace_code"`
		Astrology         string `json:"astrology"`
		Ethnicity         string `json:"ethnicity"`
		Nationality       string `json:"nationality"`
		HairColour        string `json:"hair_colour"`
		EyeColour         string `json:"eye_colour"`
		Weight            string `json:"weight"`
		Height            string `json:"height"`
		Measurements      string `json:"measurements"`
		Cupsize           string `json:"cupsize"`
		Tattoos           string `json:"tattoos"`
		Piercings         string `json:"piercings"`
		Waist             string `json:"waist"`
		Hips              string `json:"hips"`
		FakeBoobs         bool   `json:"fake_boobs"`
		SameSexOnly       bool   `json:"same_sex_only"`
		CareerStartYear   int    `json:"career_start_year"`
		CareerEndYear     int    `json:"career_end_year"`
	} `json:"extras"`
	Aliases   []string `json:"aliases"`
	Image     string   `json:"image"`
	Thumbnail string   `json:"thumbnail"`
	Face      string   `json:"face"`
	Posters   []struct {
		ID    int    `json:"id"`
		URL   string `json:"url"`
		Size  int    `json:"size"`
		Order int    `json:"order"`
	} `json:"posters"`
}

type tpdbVideoInfo struct {
	ID         string `json:"id"`
	NumericID  uint64 `json:"_id"`
	ExternalID string `json:"external_id"`
	// Slug is the meaningful id.
	Slug string `json:"slug"`

	Title       string  `json:"title"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Rating      float32 `json:"rating"`
	SiteID      int     `json:"site_id"`
	Date        string  `json:"date"`
	URL         string  `json:"url"`

	// Cover
	Image       string `json:"image"`
	BackImage   string `json:"back_image"`
	PosterImage string `json:"poster_image"`
	// Thumbnail
	Poster   string `json:"poster"`
	Trailer  string `json:"trailer"`
	Duration int    `json:"duration"`

	Performers []tpdbActorInfo `json:"performers"`
	Site       struct {
		Name string `json:"name"`
	} `json:"site"`
	Tags []struct {
		Name string `json:"name"`
	} `json:"tags"`
	Directors []struct {
		Name string `json:"name"`
	} `json:"directors"`
}

type searchTPDBVideosResponse struct {
	Data []tpdbVideoInfo `json:"data"`
}

type TPDBSearchVideosInput struct {
	Query string `json:"query" jsonschema:"the name of the video to search for, don't include release date and studio prefix, don't use dash or dot spliter"`
}

type TPDBVideoItem struct {
	ID          string   `json:"id" jsonschema:"the meaningful id of the video, usually used to rename files."`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Type        string   `json:"type" jsonschema:"scene or movie"`
	Date        string   `json:"date"`
	Actors      []string `json:"actors"`
}
type TPDBSearchVideosOutput struct {
	Results []TPDBVideoItem `json:"results"`
}

func (s *ThePornDB) search(query string, url_ string) ([]TPDBVideoItem, error) {
	u, err := url.Parse(url_)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("q", query)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+s.apiToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := searchTPDBVideosResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	var results []TPDBVideoItem
	for _, item := range res.Data {
		if len(results) >= tpdbLimitVideoPerType {
			break
		}

		var actors []string
		for _, actor := range item.Performers {
			actors = append(actors, actor.Name)
		}
		results = append(results, TPDBVideoItem{
			ID:          item.Slug,
			Title:       item.Title,
			Description: item.Description,
			Type:        item.Type,
			Date:        item.Date,
			Actors:      actors,
		})
	}
	return results, nil
}

// searchTPDBVideos will search both on scene and movie
func (s *ThePornDB) searchTPDBVideos(input TPDBSearchVideosInput) (TPDBSearchVideosOutput, error) {
	// search scene
	res, err := s.search(input.Query, tpdbSearchSceneURL)
	if err != nil {
		return TPDBSearchVideosOutput{}, err
	}
	// search movie
	searchMoviesRes, err := s.search(input.Query, tpdbSearchMovieURL)
	if err != nil {
		return TPDBSearchVideosOutput{}, err
	}
	res = append(res, searchMoviesRes...)
	return TPDBSearchVideosOutput{Results: res}, nil
}

func (s *ThePornDB) searchTPDBVideosTool(ctx context.Context, req *mcp.CallToolRequest, input TPDBSearchVideosInput) (
	*mcp.CallToolResult, TPDBSearchVideosOutput, error) {
	result, err := s.searchTPDBVideos(input)
	return nil, result, err
}
