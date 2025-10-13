package mcptools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWikipedia_searchWikipedia(t *testing.T) {
	w := NewWikipedia("en")
	input := WikipediaSearchInput{Query: "Go programming language"}
	output, err := w.searchWikipedia(input)

	assert.NoError(t, err)
	assert.NotEmpty(t, output.Results, "Expected search results to not be empty")
	assert.GreaterOrEqual(t, len(output.Results), 1, "Expected at least one search result")

	// Check if the first result has a title and summary
	if len(output.Results) > 0 {
		assert.NotEmpty(t, output.Results[0].Title, "Expected first result to have a title")
		assert.NotEmpty(t, output.Results[0].Summary, "Expected first result to have a summary")
	}
}

func TestWikipedia_wikipediaPage(t *testing.T) {
	w := NewWikipedia("en")
	input := WikipediaPageInput{Title: "Go (programming language)"}
	output, err := w.wikipediaPage(input)

	assert.NoError(t, err)
	assert.NotEmpty(t, output.Content, "Expected page content to not be empty")
}
