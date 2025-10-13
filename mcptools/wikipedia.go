package mcptools

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gowiki "github.com/trietmn/go-wiki"
)

type Wikipedia struct {
}

func NewWikipedia(language string) *Wikipedia {
	gowiki.SetUserAgent(userAgent)
	gowiki.SetLanguage(language)
	return &Wikipedia{}
}

func (w *Wikipedia) AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "wikipedia_search",
		Description: "Search for Wikipedia pages by given query.",
	}, w.searchWikipediaTool)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "wikipedia_page",
		Description: "Get the content of a Wikipedia page by given title.",
	}, w.wikipediaPageTool)
}

type WikipediaSearchInput struct {
	Query string `json:"query"`
}

type WikipediaSearchItem struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type WikipediaSearchOutput struct {
	Results []WikipediaSearchItem `json:"results"`
}

func (w *Wikipedia) searchWikipedia(input WikipediaSearchInput) (WikipediaSearchOutput, error) {
	searchResults, _, err := gowiki.Search(input.Query, 3, false)
	if err != nil {
		return WikipediaSearchOutput{}, err
	}

	var results []WikipediaSearchItem
	for _, result := range searchResults {
		item := WikipediaSearchItem{
			Title: result,
		}

		summary, err := gowiki.Summary(result, 5, -1, false, true)
		if err != nil {
			log.Printf("Error getting page %v summary: %v", result, err)
		} else {
			item.Summary = summary
		}
		results = append(results, item)
	}

	return WikipediaSearchOutput{Results: results}, nil
}

func (w *Wikipedia) searchWikipediaTool(
	ctx context.Context, req *mcp.CallToolRequest, input WikipediaSearchInput) (
	*mcp.CallToolResult, WikipediaSearchOutput, error) {
	result, err := w.searchWikipedia(input)
	return nil, result, err
}

type WikipediaPageInput struct {
	Title string `json:"title"`
}

type WikipediaPageOutput struct {
	Content string `json:"content"`
}

func (w *Wikipedia) wikipediaPage(input WikipediaPageInput) (WikipediaPageOutput, error) {
	page, err := gowiki.GetPage(input.Title, -1, true, true)
	if err != nil {
		return WikipediaPageOutput{}, err
	}

	content, err := page.GetContent()
	if err != nil {
		return WikipediaPageOutput{}, err
	}

	return WikipediaPageOutput{Content: content}, nil
}

func (w *Wikipedia) wikipediaPageTool(
	ctx context.Context, req *mcp.CallToolRequest, input WikipediaPageInput) (
	*mcp.CallToolResult, WikipediaPageOutput, error) {
	result, err := w.wikipediaPage(input)
	return nil, result, err
}
