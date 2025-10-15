package mcptools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchAlias(t *testing.T) {
	alias := NewJAVActorAlias("")
	input := JAVActorAliasInput{Name: "藤森里穂"}
	output, err := alias.searchAlias(t.Context(), input)

	require.NoError(t, err)
	assert.Contains(t, output.Aliases, "藤森里穂")
	assert.Contains(t, output.Aliases, "井上遥香")
}

func TestNameToDir(t *testing.T) {
	// Create a temporary JSON file for testing
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_aliases.json")

	initialData := map[string][]string{
		"ActorA": {"Alias1", "Alias2"},
		"ActorB": {"Alias3", "Alias4"},
	}
	dataBytes, err := json.MarshalIndent(initialData, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(tempFile, dataBytes, 0644)
	require.NoError(t, err)

	alias := NewJAVActorAlias(tempFile)

	// Test case 1: Known alias
	input1 := NameToDir{Name: "Alias1"}
	output1, err := alias.nameToDir(input1)
	require.NoError(t, err)
	assert.Equal(t, "ActorA", output1.Dir)

	// Test case 2: Unknown alias
	input2 := NameToDir{Name: "UnknownAlias"}
	output2, err := alias.nameToDir(input2)
	require.NoError(t, err)
	assert.Empty(t, output2.Dir)
}

func TestAddAlias(t *testing.T) {
	// Create a temporary JSON file for testing
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_aliases.json")

	initialData := map[string][]string{
		"ActorA": {"Alias1", "Alias2"},
	}
	dataBytes, err := json.MarshalIndent(initialData, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(tempFile, dataBytes, 0644)
	require.NoError(t, err)

	alias := NewJAVActorAlias(tempFile)

	// Test case 1: Add aliases to an existing actor
	input1 := AddAliasInput{
		Name:    "ActorA",
		Aliases: []string{"Alias1", "Alias2", "NewAlias3"},
	}
	output1, err := alias.addAlias(input1)
	require.NoError(t, err)
	assert.Equal(t, "ActorA", output1.Dir)

	// Verify the file content
	updatedData1, _, err := alias.loadJSONFile()
	require.NoError(t, err)
	assert.Contains(t, updatedData1["ActorA"], "NewAlias3")
	assert.Len(t, updatedData1["ActorA"], 3) // Alias1, Alias2, NewAlias3

	// Test case 2: Add a new actor
	input2 := AddAliasInput{
		Name:    "ActorC",
		Aliases: []string{"Alias5", "Alias6"},
	}
	output2, err := alias.addAlias(input2)
	require.NoError(t, err)
	assert.Equal(t, "ActorC", output2.Dir)

	// Verify the file content
	updatedData2, _, err := alias.loadJSONFile()
	require.NoError(t, err)
	assert.Contains(t, updatedData2, "ActorC")
	assert.Contains(t, updatedData2["ActorC"], "Alias5")
	assert.Contains(t, updatedData2["ActorC"], "Alias6")
	assert.Len(t, updatedData2["ActorC"], 2)
}
