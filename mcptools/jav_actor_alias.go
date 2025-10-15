package mcptools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type JAVActorAlias struct {
	jsonFile string
}

func NewJAVActorAlias(jsonFile string) *JAVActorAlias {
	return &JAVActorAlias{jsonFile: jsonFile}
}

func (s *JAVActorAlias) AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "web_search_jav_actor_alias",
		Description: "Searches JAVDB for an actor's aliases given one of their names.",
	}, s.searchAliasTool)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "jav_actor_name_to_dir",
		Description: "Determines the directory name used for a JAV actor, given one of their alias names.",
	}, s.nameToDirTool)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "jav_actor_add_alias",
		Description: "Adds or updates JAV actor aliases in the system and returns the actor's directory name.",
	}, s.addAliasTool)
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

type NameToDir struct {
	Name string `json:"name" jsonschema:"the name of the actor to search for"`
}

type NameToDirOutput struct {
	Dir string `json:"dir,omitempty" jsonschema:"the dir of the actor given, maybe empty if not found"`
}

// loadJSONFile return dir -> alias mapping and alias -> dir mapping.
func (s *JAVActorAlias) loadJSONFile() (map[string][]string, map[string]string, error) {
	d, err := os.ReadFile(s.jsonFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}
	jsonFile := map[string][]string{}
	if err := json.Unmarshal(d, &jsonFile); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}
	aliasToDir := map[string]string{}
	for dir, aliases := range jsonFile {
		for _, alias := range aliases {
			aliasToDir[alias] = dir
		}
	}
	return jsonFile, aliasToDir, nil
}

func (s *JAVActorAlias) nameToDir(input NameToDir) (NameToDirOutput, error) {
	_, aliasToDir, err := s.loadJSONFile()
	if err != nil {
		return NameToDirOutput{}, err
	}
	dir, ok := aliasToDir[input.Name]
	if !ok {
		return NameToDirOutput{}, nil
	}
	return NameToDirOutput{Dir: dir}, nil
}

func (s *JAVActorAlias) nameToDirTool(
	ctx context.Context, req *mcp.CallToolRequest, input NameToDir) (
	*mcp.CallToolResult, NameToDirOutput, error) {
	output, err := s.nameToDir(input)
	return &mcp.CallToolResult{}, output, err
}

type AddAliasInput struct {
	Name    string   `json:"name" jsonschema:"the best-known name of the actor, this maybe use as dir name"`
	Aliases []string `json:"aliases" jsonschema:"all aliases of the actor, should also include the best-known name"`
}

type AddAliasOutput struct {
	Dir string `json:"dir" jsonschema:"the dir of the actor given"`
}

func (s *JAVActorAlias) addAlias(input AddAliasInput) (AddAliasOutput, error) {
	dirToAlias, aliasToDir, err := s.loadJSONFile()
	if err != nil {
		return AddAliasOutput{}, err
	}

	dir := ""
	for _, alias := range input.Aliases {
		d, ok := aliasToDir[alias]
		if ok {
			dir = d
			break
		}
	}

	if dir != "" {
		// if dir is found, actor got new name, we need to append new names to this actor.
		aliases := dirToAlias[dir]
		for _, a := range input.Aliases {
			if !slices.Contains(aliases, a) {
				aliases = append(aliases, a)
			}
		}
		dirToAlias[dir] = aliases
	} else {
		// if dir not found, need to add new actor.
		dir = input.Name
		dirToAlias[dir] = input.Aliases
	}

	out, err := json.MarshalIndent(dirToAlias, "", "  ")
	if err != nil {
		return AddAliasOutput{}, fmt.Errorf("failed to marshal json: %w", err)
	}

	if err := os.WriteFile(s.jsonFile, out, 0644); err != nil {
		return AddAliasOutput{}, fmt.Errorf("failed to write file: %w", err)
	}

	return AddAliasOutput{Dir: dir}, nil
}

func (s *JAVActorAlias) addAliasTool(
	ctx context.Context, req *mcp.CallToolRequest, input AddAliasInput) (
	*mcp.CallToolResult, AddAliasOutput, error) {
	output, err := s.addAlias(input)
	return &mcp.CallToolResult{}, output, err
}
