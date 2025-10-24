package mcptools

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	avabse = "AVBASE"
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
		Description: "Searches for Japanese and Chinese pornographic content on Metatube using a given ID (番号), e.g., 'SSIS-698'.",
	}, s.searchJAVTool)
}

type SearchJAVInput struct {
	JAVID string `json:"jav_id" jsonschema:"the id (番号) of the jav to search for, it usually Studio/Label Prefix (usually 3-4 letters) then dash (-) then number. for example: SSIS-698"`
}

type JAV struct {
	JAVID       string   `json:"jav_id"`
	Title       string   `json:"title"`
	Provider    string   `json:"provider"`
	Actors      []string `json:"actors,omitempty"`
	ReleaseDate string   `json:"release_date"`
	Tags        []string `json:"tags,omitempty"`
	Maker       string   `json:"maker,omitempty"`
	Label       string   `json:"label,omitempty"`
	Series      string   `json:"series,omitempty"`
}

type SearchJAVOutput struct {
	Results []JAV `json:"results"`
}

type MetatubeJAVSearchResponse struct {
	Data []struct {
		ID          string   `json:"id"`
		Number      string   `json:"number"`
		Title       string   `json:"title"`
		Provider    string   `json:"provider"`
		Actors      []string `json:"actors,omitempty"`
		ReleaseDate string   `json:"release_date"`
	} `json:"data"`
}

type MetatubeJAVDetaiisResponse struct {
	Data struct {
		Maker  string   `json:"maker,omitempty"`
		Label  string   `json:"label,omitempty"`
		Series string   `json:"series,omitempty"`
		Genres []string `json:"genres,omitempty"`
	} `json:"data"`
}

func (s *Metatube) searchJAV(ctx context.Context, input SearchJAVInput) (SearchJAVOutput, error) {
	u, err := url.Parse(s.apiURL)
	if err != nil {
		return SearchJAVOutput{}, err
	}
	u.Path = "/v1/movies/search"
	u.RawQuery = url.Values{
		"q": {input.JAVID},
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

	res := MetatubeJAVSearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return SearchJAVOutput{}, err
	}
	var results []JAV
	for _, item := range res.Data {
		jav := JAV{
			JAVID:       item.Number,
			Title:       item.Title,
			Provider:    item.Provider,
			Actors:      item.Actors,
			ReleaseDate: item.ReleaseDate,
		}

		// AVBASE usually has best metadata, pull details for tags.
		if item.Provider == avabse {
			u, err := url.Parse(s.apiURL)
			if err != nil {
				log.Printf("Error parsing URL: %v", err)
				continue
			}
			u.Path = path.Join("/v1/movies/", avabse, item.ID)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
			if err != nil {
				log.Printf("Error creating request: %v", err)
				continue
			}
			if s.apiKey != "" {
				req.Header.Add("Authorization", "Bearer "+s.apiKey)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return SearchJAVOutput{}, err
			}
			defer resp.Body.Close()

			detailsRes := MetatubeJAVDetaiisResponse{}
			err = json.NewDecoder(resp.Body).Decode(&detailsRes)
			if err != nil {
				log.Printf("Error decoding response: %v", err)
				continue
			}
			jav.Maker = detailsRes.Data.Maker
			jav.Label = detailsRes.Data.Label
			jav.Series = detailsRes.Data.Series
			jav.Tags = detailsRes.Data.Genres
		}

		results = append(results, jav)
	}
	return SearchJAVOutput{Results: results}, nil
}

func (s *Metatube) searchJAVTool(
	ctx context.Context, req *mcp.CallToolRequest, input SearchJAVInput) (
	*mcp.CallToolResult, SearchJAVOutput, error) {
	result, err := s.searchJAV(ctx, input)
	return nil, result, err
}
