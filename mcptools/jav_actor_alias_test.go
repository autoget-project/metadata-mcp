package mcptools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchAlias(t *testing.T) {
	alias := NewJAVActorAlias()
	input := JAVActorAliasInput{Name: "藤森里穂"}
	output, err := alias.searchAlias(t.Context(), input)

	assert.NoError(t, err)
	assert.Contains(t, output.Aliases, "藤森里穂")
	assert.Contains(t, output.Aliases, "井上遥香")
}
