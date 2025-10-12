package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/autoget-project/metadata-mcp/config"
	"github.com/autoget-project/metadata-mcp/mcptools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	configPath := flag.String("c", "", "path to config file")
	flag.Parse()

	var conf *config.Config
	var err error

	if *configPath == "" {
		conf, err = config.ReadConfig(*configPath)
		if err != nil {
			log.Fatalf("Error reading config from %s: %v", *configPath, err)
		}
	} else {
		conf, err = config.ReadConfigFromEnv()
		if err != nil {
			log.Fatalf("Error reading config from environment variables: %v", err)
		}
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name: "metadata-mcp-server",
	}, nil)

	mcptools.NewTMDB(conf.TMDBAPIKey, conf.TMDBResponseLanguage).AddTools(server)
	mcptools.NewThePornDB(conf.ThePornDBAPIToken).AddTools(server)
	mcptools.NewMetatube(conf.MetaTubeAPIURL, conf.MetaTubeAPIKEY).AddTools(server)
	// TODO: Add other MCP tools here (ThePornDB, Metatube, DuckDuckGo)

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", conf.Port),
		Handler: handler,
	}

	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on port %v", conf.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-stop // Wait for OS signal

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
