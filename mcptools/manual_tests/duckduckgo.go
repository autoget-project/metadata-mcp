package main

import (
	"context"

	"github.com/autoget-project/metadata-mcp/mcptools"
)

func main() {
	ddg, err := mcptools.NewDuckDuckGo()
	if err != nil {
		panic(err)
	}
	_, res, err := ddg.SearchDuckDuckGoTool(context.Background(), nil, mcptools.DuckDuckGoSearchInput{Query: "三上悠亚"})
	if err != nil {
		panic(err)
	}
	println(res.Results)
}
