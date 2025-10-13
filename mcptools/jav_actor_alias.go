package mcptools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type JAVActorAlias struct {
}

func NewJAVActorAlias() *JAVActorAlias {
	return &JAVActorAlias{}
}

func (s *JAVActorAlias) AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "web_search_jav_actor_alias",
		Description: "JAV actor usually has many alias name, use one of their name find alias name on JAVDB",
	}, s.searchAliasTool)
	// TODO: get local dir -> actor alias
	// TODO: update local dir -> actor alias
}

type JAVActorAliasInput struct {
	Name string `json:"name" jsonschema:"the name of the actor to search for"`
}

type JAVActorAliasOutput struct {
	Aliases []string `json:"aliases"`
}

func (s *JAVActorAlias) httpGet(ctx context.Context, url_ string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url_, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL, status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	content := string(bodyBytes)
	return content, nil
}

// searchAlias get actor alias in following steps:
// GET https://javdb.com/search?f=actor&q=$NAME, find the actor alias in actor-box
func (s *JAVActorAlias) searchAlias(ctx context.Context, input JAVActorAliasInput) (JAVActorAliasOutput, error) {
	u, err := url.Parse("https://javdb.com/search?f=actor")
	if err != nil {
		return JAVActorAliasOutput{}, fmt.Errorf("failed to parse URL: %w", err)
	}
	q := u.Query()
	q.Set("q", input.Name)
	u.RawQuery = q.Encode()
	content, err := s.httpGet(ctx, u.String())
	if err != nil {
		return JAVActorAliasOutput{}, fmt.Errorf("failed to fetch URL: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return JAVActorAliasOutput{}, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var aliases []string
	doc.Find(".actor-box a").First().Each(func(i int, s *goquery.Selection) {
		title, exists := s.Attr("title")
		if exists {
			aliases = strings.Split(title, ", ")
		}
	})

	return JAVActorAliasOutput{Aliases: aliases}, nil
}

func (s *JAVActorAlias) searchAliasTool(
	ctx context.Context, req *mcp.CallToolRequest, input JAVActorAliasInput) (
	*mcp.CallToolResult, JAVActorAliasOutput, error) {
	output, err := s.searchAlias(ctx, input)
	if err != nil {
		return nil, JAVActorAliasOutput{}, err
	}
	return &mcp.CallToolResult{}, output, nil
}
