package mcptools

import (
	"context"
	"fmt"
	"io"
	"net/http"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Fetcher struct {
}

func NewFetcher() *Fetcher {
	return &Fetcher{}
}

func (f *Fetcher) AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "fetch",
		Description: "fetch content from given url",
	}, f.fetchTool)
}

type FetchInput struct {
	URL               string `json:"url" jsonschema:"the url to fetch"`
	ConvertToMarkdown bool   `json:"convert_to_markdown" jsonschema:"(optional) whether to convert the content to markdown, default is no"`
}

type FetchOutput struct {
	Content string `json:"content"`
}

func (f *Fetcher) fetch(ctx context.Context, input FetchInput) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, input.URL, nil)
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

	if input.ConvertToMarkdown {
		markdown, err := htmltomarkdown.ConvertString(content)
		if err != nil {
			return "", fmt.Errorf("failed to convert HTML to markdown: %w", err)
		}
		content = markdown
	}

	return content, nil
}

func (f *Fetcher) fetchTool(ctx context.Context, req *mcp.CallToolRequest, input FetchInput) (
	*mcp.CallToolResult, FetchOutput, error) {
	result, err := f.fetch(ctx, input)
	return nil, FetchOutput{Content: result}, err
}
