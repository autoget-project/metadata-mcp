package mcptools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestFetcher_fetch(t *testing.T) {
	t.Run("successful fetch", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Hello, World!"))
		}))
		defer server.Close()

		fetcher := &Fetcher{}
		input := FetchInput{
			URL: server.URL,
		}

		content, err := fetcher.fetch(t.Context(), input)
		assert.NoError(t, err)
		assert.Equal(t, "Hello, World!", content)
	})

	t.Run("successful fetch and convert to markdown", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("<h1>Title</h1><p>Content</p>"))
		}))
		defer server.Close()

		fetcher := &Fetcher{}
		input := FetchInput{
			URL:               server.URL,
			ConvertToMarkdown: true,
		}

		content, err := fetcher.fetch(t.Context(), input)
		assert.NoError(t, err)
		assert.Contains(t, content, "# Title")
		assert.Contains(t, content, "Content")
	})

	t.Run("failed fetch - bad URL", func(t *testing.T) {
		fetcher := &Fetcher{}
		input := FetchInput{
			URL: "http://localhost:99999", // Non-existent URL
		}

		content, err := fetcher.fetch(t.Context(), input)
		assert.Error(t, err)
		assert.Empty(t, content)
		assert.Contains(t, err.Error(), "failed to fetch URL")
	})

	t.Run("failed fetch - non-200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		fetcher := &Fetcher{}
		input := FetchInput{
			URL: server.URL,
		}

		content, err := fetcher.fetch(t.Context(), input)
		assert.Error(t, err)
		assert.Empty(t, content)
		assert.Contains(t, err.Error(), "failed to fetch URL, status code: 404")
	})
}

func TestFetcher_fetchTool(t *testing.T) {
	t.Run("successful tool call", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Tool Test Content"))
		}))
		defer server.Close()

		fetcher := &Fetcher{}
		input := FetchInput{
			URL: server.URL,
		}

		ctx := context.Background()
		req := &mcp.CallToolRequest{}

		result, output, err := fetcher.fetchTool(ctx, req, input)
		assert.NoError(t, err)
		assert.Nil(t, result) // User removed mcp.CallToolSuccess()
		assert.Equal(t, "Tool Test Content", output.Content)
	})

	t.Run("failed tool call - fetch error", func(t *testing.T) {
		fetcher := &Fetcher{}
		input := FetchInput{
			URL: "http://localhost:99999", // Non-existent URL
		}

		ctx := context.Background()
		req := &mcp.CallToolRequest{}

		result, output, err := fetcher.fetchTool(ctx, req, input)
		assert.Error(t, err)
		assert.Empty(t, output.Content)
		assert.Nil(t, result) // User removed mcp.CallToolError(err)
		assert.Contains(t, err.Error(), "failed to fetch URL")
	})
}
