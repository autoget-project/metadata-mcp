package mcptools

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func metadataURLFromEnv(t *testing.T) string {
	t.Helper()
	url := os.Getenv("METADATA_URL")
	if url == "" {
		t.Skip("METADATA_URL environment variable not set")
	}
	return url
}

func TestSearchJAV(t *testing.T) {
	url := metadataURLFromEnv(t)
	metatube := NewMetatube(url, "")
	result, err := metatube.searchJAV(SearchJAVInput{JAVID: "SSIS-698"})
	require.NoError(t, err)
	require.NotEmpty(t, result.Results)
	assert.Equal(t, "SSIS-698", result.Results[0].JAVID)
	assert.Contains(t, result.Results[0].Actors, "三上悠亜")
}
