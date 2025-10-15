package mcptools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tmc/langchaingo/tools/duckduckgo"
)

const (
	userAgent          = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36"
	ddgMaxSearchResult = 10
)

type DuckDuckGo struct {
	tool *duckduckgo.Tool
}

func NewDuckDuckGo() (*DuckDuckGo, error) {
	tool, err := duckduckgo.New(ddgMaxSearchResult, userAgent)
	if err != nil {
		return nil, err
	}
	return &DuckDuckGo{
		tool: tool,
	}, nil
}

func (s *DuckDuckGo) AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "web_search",
		Description: "Performs a web search using DuckDuckGo and returns the search results.",
	}, s.SearchDuckDuckGoTool)
}

type DuckDuckGoSearchInput struct {
	Query string `json:"query"`
}

type DuckDuckGoSearchOutput struct {
	Results string `json:"results"`
}

func (s *DuckDuckGo) SearchDuckDuckGoTool(
	ctx context.Context, req *mcp.CallToolRequest, input DuckDuckGoSearchInput) (
	*mcp.CallToolResult, DuckDuckGoSearchOutput, error) {
	res, err := s.tool.Call(ctx, input.Query)
	return nil, DuckDuckGoSearchOutput{Results: res}, err
}
