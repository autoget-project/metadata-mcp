# Metadata MCP Server

A metadata server that searches for metadata from various sources, including TMDB, ThePornDB, and Metatube.

## Features

*   **Comprehensive Movie & TV Show Search:** Utilizes The Movie Database (TMDB) to find detailed metadata for movies and TV shows, including actors, release dates, and overviews.
*   **Specialized Pornographic Metadata:** Integrates with ThePornDB for extensive search capabilities for non-Japanese pornographic content.
*   **JAV Content Discovery:** Connects to Metatube for specialized search and metadata retrieval for Japanese Adult Video (JAV) content.
*   **JAV Actor Alias Management:** Provides tools to search for, resolve, and manage aliases for JAV actors, ensuring consistent naming and organization.
*   **General Web Search Fallback:** Includes DuckDuckGo for broader web searches when specialized metadata sources may not cover a query.
*   **Wikipedia Integration:** Offers tools to search for and retrieve content from Wikipedia pages for general information.
*   **URL Content Fetching:** Allows fetching content from any given URL, with an option to convert HTML to Markdown for easier readability.

## Installation

```
docker run \
  -p 8080:8080 \
  -v "/path/to/jav_actor_aliases.json:/app/jav_actor_aliases.json" \
  -e TMDB_API_KEY="<YOUR_TMDB_API_KEY>" \
  -e TPDB_API_TOKEN="<YOUR_TPDB_API_TOKEN>" \
  -e METATUBE_API_URL="<YOUR_METATUBE_API_URL>" \
  -e JAV_ACTOR_JSON_FILEPATH="/app/jav_actor_aliases.json" \
  ghcr.io/autoget-project/metadata-mcp:main
```

### Config: Environment Variables

This server uses environment variables for configuration. The following variables are available:

*   `PORT` (optional): The port the server will listen on. Defaults to `8080`.
*   `TMDB_API_KEY` (required): Your API key for The Movie Database (TMDB).
*   `TMDB_RESPONSE_LANGUAGE` (optional): The language for TMDB responses. Defaults to `zh-CN`.
*   `TPDB_API_TOKEN` (required): Your API token for ThePornDB.
*   `METATUBE_API_URL` (required): The base URL for the Metatube API.
*   `METATUBE_API_KEY` (optional): Your API key for Metatube.
*   `WIKIPEDIA_LANGUAGE` (optional): The language for Wikipedia searches. Defaults to `zh`.
*   `JAV_ACTOR_JSON_FILEPATH` (required): The path to the JSON file used for storing JAV actor aliases.

## Tools

The Metadata MCP Server exposes the following tools:

*   **web_search**: Performs a web search using DuckDuckGo and returns the search results.
*   **fetch**: Fetches content from a specified URL. Can optionally convert HTML content to Markdown.
*   **web_search_jav_actor_alias**: Searches JAVDB to find all known aliases for a given JAV actor's name.
*   **jav_actor_name_to_dir**: Retrieves the canonical directory name for a JAV actor based on one of their alias names.
*   **jav_actor_add_alias**: Manages JAV actor aliases by adding new ones or updating existing entries, returning the actor's primary directory name.
*   **search_japanese_porn**: Searches for Japanese and Chinese pornographic content on Metatube using a given ID (番号), e.g., 'SSIS-698'.
*   **search_porn**: Searches for non-Japanese pornographic movies and scenes on ThePornDB.
*   **search_movies**: Searches for movies on The Movie Database (TMDB) by name (required) and optional release year.
*   **search_tv_shows**: Searches for TV shows on The Movie Database (TMDB) by name.
*   **wikipedia_search**: Searches Wikipedia for pages matching a given query and returns a summary of each result.
*   **wikipedia_page**: Retrieves the full content of a Wikipedia page given its exact title.
