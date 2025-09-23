// Package main ...
package main

import (
	"fmt"
	"log"
	"os"

	"immich-fave-to-album/immich"
)

func main() {
	config := parseFlags()

	// Get API key from environment
	config.APIKey = os.Getenv("IMMICH_API_KEY")
	if config.APIKey == "" {
		log.Fatal("IMMICH_API_KEY environment variable is required")
	}

	// Validate configuration
	if err := immich.ValidateConfig(config); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	if config.Verbose {
		fmt.Printf("Verbose mode enabled\n")
		fmt.Printf("Dry run mode: %t\n", config.DryRun)
		fmt.Printf("Connecting to Immich server: %s\n", config.ServerURL)
		fmt.Printf("Album name: %s\n", config.AlbumName)
	}

	if config.Debug {
		fmt.Printf("Debug mode enabled\n")
	}

	// Create Immich client
	client := immich.NewClient(config)

	// Get all favorite assets
	assets, err := client.GetFavoriteAssets()
	if err != nil {
		log.Fatalf("Failed to get favorite assets: %v", err)
	}

	if len(assets) == 0 {
		fmt.Println("No favorite assets found")
		return
	}

	if config.Verbose {
		fmt.Printf("Found %d favorite assets\n", len(assets))
	}

	// Create or get album
	album, err := client.CreateOrGetAlbum(assets)
	if err != nil {
		log.Fatalf("Failed to create/get album: %v", err)
	}

	if config.Verbose {
		fmt.Printf("Using album: %s (ID: %s)\n", album.AlbumName, album.ID)
	}

	// Add assets to album if not already added during creation
	if len(assets) > 0 {
		if config.DryRun {
			fmt.Printf("[DRY RUN] Would add %d assets to album '%s'\n", len(assets), album.AlbumName)
			for _, asset := range assets {
				fmt.Printf("[DRY RUN] Would add asset: %s (type: %s, favorite: %t)\n", asset.ID, asset.Type, asset.IsFavorite)
			}
			// Always show final status for dry run
			fmt.Printf("✓ Dry run completed - would add %d assets to album '%s'\n", len(assets), album.AlbumName)
		} else {
			err = client.AddAssetsToAlbum(album.ID, assets)
			if err != nil {
				log.Fatalf("Failed to add assets to album: %v", err)
			}
			// Always show final success status
			fmt.Printf("✓ Successfully added %d assets to album '%s'\n", len(assets), album.AlbumName)
		}
	}
}

func parseFlags() immich.Config {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go -server <immich-server-url> -album <album-name> [options]")
		fmt.Println("Options:")
		fmt.Println("  -verbose       Enable verbose output (default: false)")
		fmt.Println("  -debug         Enable debug output with HTTP request/response details (default: false)")
		fmt.Println("  -no-dry-run    Disable dry run mode (default: enabled)")
		fmt.Println("Example: go run main.go -server http://localhost:2283 -album 'My Favorites' -verbose -debug")
		fmt.Println("Environment variable IMMICH_API_KEY must be set")
		os.Exit(1)
	}

	config := immich.Config{
		DryRun: true, // Default to dry run mode
	}
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-server":
			if i+1 < len(os.Args) {
				config.ServerURL = os.Args[i+1]
				i++
			}
		case "-album":
			if i+1 < len(os.Args) {
				config.AlbumName = os.Args[i+1]
				i++
			}
		case "-verbose":
			config.Verbose = true
		case "-debug":
			config.Debug = true
		case "-no-dry-run":
			config.DryRun = false
		}
	}

	if config.ServerURL == "" || config.AlbumName == "" {
		fmt.Println("Both -server and -album flags are required")
		os.Exit(1)
	}

	return config
}
