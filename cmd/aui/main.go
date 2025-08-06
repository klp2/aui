package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/aui/internal/config"
	"github.com/yourusername/aui/internal/storage"
	"github.com/yourusername/aui/internal/ui"
)

func main() {
	// Define command-line flags
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	var cfg *config.Config
	var err error

	if configPath != "" {
		// Load from specified config file
		cfg, err = config.LoadFromFile(configPath)
		if err != nil {
			log.Fatalf("Failed to load config from %s: %v", configPath, err)
		}
		// Apply environment overrides even when using custom config
		cfg.LoadFromEnv()
	} else {
		// Load from default location or create default config
		cfg, err = config.Load()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Ensure database directory exists
	dbDir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Initialize storage
	store, err := storage.NewSQLiteStore(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()

	// Create the initial app state with config and storage
	app := ui.InitialAppWithDependencies(cfg, store)

	// Set up logging if configured
	if cfg.Logging.File != "" {
		logDir := filepath.Dir(cfg.Logging.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Printf("Warning: Failed to create log directory: %v", err)
		} else {
			logFile, err := os.OpenFile(cfg.Logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				log.Printf("Warning: Failed to open log file: %v", err)
			} else {
				defer logFile.Close()
				log.SetOutput(logFile)
			}
		}
	}

	// Create the Bubble Tea program
	p := tea.NewProgram(app, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running aui: %v\n", err)
		os.Exit(1)
	}
}
