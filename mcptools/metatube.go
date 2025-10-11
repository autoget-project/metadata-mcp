package mcptools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Metatube struct {
	apiURL string
	apiKey string
}

func NewMetatube(apiURL, apiKey string) *Metatube {
	return &Metatube{
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

func (s *Metatube) AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_japanese_porn",
		Description: "Search for Japanese pornographic by given JavID (番号).",
	}, s.searchTJAVTool)
}

type SearchJAVInput struct {
	JAVID string `json:"jav_id" jsonschema:"the id (番号) of the jav to search for, it usually Studio/Label Prefix (usually 3-4 letters) then dash (-) then number. for example: SSIS-698"`
}

type JAV struct {
	JAVID       string   `json:"jav_id"`
	Title       string   `json:"title"`
	Provider    string   `json:"provider"`
	Actors      []string `json:"actors"`
	ReleaseDate string   `json:"release_date"`
}

type SearchJAVOutput struct {
	Results []JAV `json:"results"`
}

type MetatubeSearchJAVResponse struct {
	Data []struct {
		ID          string   `json:"id"`
		Number      string   `json:"number"`
		Title       string   `json:"title"`
		Provider    string   `json:"provider"`
		Homepage    string   `json:"homepage"`
		ThumbURL    string   `json:"thumb_url"`
		CoverURL    string   `json:"cover_url"`
		Score       int      `json:"score"`
		Actors      []string `json:"actors,omitempty"`
		ReleaseDate string   `json:"release_date"`
	} `json:"data"`
}

func (s *Metatube) searchJAV(input SearchJAVInput) (SearchJAVOutput, error) {
	u, err := url.Parse(s.apiURL)
	if err != nil {
		return SearchJAVOutput{}, err
	}
	u.Path = "/v1/movies/search"
	u.RawQuery = url.Values{
		"q": {input.JAVID},
	}.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return SearchJAVOutput{}, err
	}

	if s.apiKey != "" {
		req.Header.Add("Authorization", "Bearer "+s.apiKey)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SearchJAVOutput{}, err
	}
	defer resp.Body.Close()

	res := MetatubeSearchJAVResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return SearchJAVOutput{}, err
	}
	var results []JAV
	for _, item := range res.Data {
		results = append(results, JAV{
			JAVID:       item.Number,
			Title:       item.Title,
			Provider:    item.Provider,
			Actors:      item.Actors,
			ReleaseDate: item.ReleaseDate,
		})
	}
	return SearchJAVOutput{Results: results}, nil
}

func (s *Metatube) searchTJAVTool(
	ctx context.Context, req *mcp.CallToolRequest, input SearchJAVInput) (
	*mcp.CallToolResult, SearchJAVOutput, error) {
	result, err := s.searchJAV(input)
	return nil, result, err
}
