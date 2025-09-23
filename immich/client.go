package immich

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client represents an Immich API client
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient creates a new Immich API client
func NewClient(config Config) *Client {
	return &Client{
		config:     config,
		httpClient: &http.Client{},
	}
}

// GetFavoriteAssets retrieves all favorite assets from the Immich server
func (c *Client) GetFavoriteAssets() ([]Asset, error) {
	if c.config.Verbose {
		fmt.Printf("Searching for favorite assets using search/metadata endpoint...\n")
	}

	// Use the search/metadata endpoint to find favorite assets
	searchURL := fmt.Sprintf("%s/api/search/metadata", c.config.ServerURL)

	searchReq := SearchMetadataRequest{
		IsFavorite: true,
		Page:       0,
		Size:       1000, // Get up to 1000 assets at once
	}

	if c.config.Verbose {
		fmt.Printf("Searching for favorites using: %s\n", searchURL)
	}

	resp, err := c.makeAPIRequest("POST", searchURL, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to search for favorite assets: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var searchResp SearchMetadataResponse
	if err := c.parseJSONResponse(resp.Body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %v", err)
	}

	if c.config.Verbose {
		fmt.Printf("Found %d favorite assets\n", len(searchResp.Assets.Items))
	}

	return searchResp.Assets.Items, nil
}

// CreateOrGetAlbum creates a new album or returns an existing one with the same name
func (c *Client) CreateOrGetAlbum(assets []Asset) (*Album, error) {
	if c.config.Verbose {
		fmt.Printf("Looking for existing album '%s'...\n", c.config.AlbumName)
	}

	// First, try to find existing album with the same name
	existingAlbum, err := c.getAlbumByName()
	if err == nil && existingAlbum != nil {
		if c.config.Verbose {
			fmt.Printf("Found existing album: %s (ID: %s)\n", existingAlbum.AlbumName, existingAlbum.ID)
		}
		return existingAlbum, nil
	}

	if c.config.Verbose {
		fmt.Printf("Album '%s' not found, will create new album\n", c.config.AlbumName)
	}

	// If not found, create a new album
	assetIDs := make([]string, len(assets))
	for i, asset := range assets {
		assetIDs[i] = asset.ID
	}

	if c.config.DryRun {
		if c.config.Verbose {
			fmt.Printf("[DRY RUN] Would create album '%s' with %d assets\n", c.config.AlbumName, len(assets))
		}
		// Return a mock album for dry run
		return &Album{
			ID:          "dry-run-album-id",
			AlbumName:   c.config.AlbumName,
			Description: fmt.Sprintf("Album created by immich-fave-to-album script on %s", time.Now().Format("2006-01-02 15:04:05")),
		}, nil
	}

	createReq := CreateAlbumRequest{
		AlbumName:   c.config.AlbumName,
		AssetIDs:    assetIDs,
		Description: fmt.Sprintf("Album created by immich-fave-to-album script on %s", time.Now().Format("2006-01-02 15:04:05")),
		AlbumUsers: []AlbumUser{
			{
				UserID: "me", // This might need to be the actual user ID
				Role:   "editor",
			},
		},
	}

	if c.config.Verbose {
		fmt.Printf("Creating album '%s'...\n", c.config.AlbumName)
	}

	albumURL := fmt.Sprintf("%s/api/albums", c.config.ServerURL)
	resp, err := c.makeAPIRequest("POST", albumURL, createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create album: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create album, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var newAlbum Album
	if err := c.parseJSONResponse(resp.Body, &newAlbum); err != nil {
		return nil, fmt.Errorf("failed to parse create album response: %v", err)
	}

	if c.config.Verbose {
		fmt.Printf("Created new album: %s (ID: %s)\n", newAlbum.AlbumName, newAlbum.ID)
	}

	return &newAlbum, nil
}

// AddAssetsToAlbum adds the given assets to the specified album
func (c *Client) AddAssetsToAlbum(albumID string, assets []Asset) error {
	if c.config.DryRun {
		if c.config.Verbose {
			fmt.Printf("[DRY RUN] Would add %d assets to album %s\n", len(assets), albumID)
		}
		return nil
	}

	assetIDs := make([]string, len(assets))
	for i, asset := range assets {
		assetIDs[i] = asset.ID
	}

	if c.config.Verbose {
		fmt.Printf("Adding %d assets to album %s...\n", len(assets), albumID)
	}

	addReq := AddAssetsRequest{
		AssetIDs: assetIDs,
	}

	addURL := fmt.Sprintf("%s/api/albums/%s/assets", c.config.ServerURL, albumID)
	resp, err := c.makeAPIRequest("PUT", addURL, addReq)
	if err != nil {
		return fmt.Errorf("failed to add assets to album: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add assets, status: %d, body: %s", resp.StatusCode, string(body))
	}

	if c.config.Verbose {
		fmt.Printf("Successfully added %d assets to album\n", len(assets))
	}

	return nil
}

// getAlbumByName retrieves an album by name
func (c *Client) getAlbumByName() (*Album, error) {
	albumsURL := fmt.Sprintf("%s/api/albums", c.config.ServerURL)
	resp, err := c.makeAPIRequest("GET", albumsURL, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get albums, status: %d", resp.StatusCode)
	}

	var albums []Album
	if err := c.parseJSONResponse(resp.Body, &albums); err != nil {
		return nil, fmt.Errorf("failed to parse albums response: %v", err)
	}

	for _, album := range albums {
		if album.AlbumName == c.config.AlbumName {
			return &album, nil
		}
	}

	return nil, fmt.Errorf("album not found")
}

// makeAPIRequest makes an HTTP request to the Immich API
func (c *Client) makeAPIRequest(method, url string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.config.APIKey)

	if c.config.Debug {
		fmt.Printf("[DEBUG] === HTTP REQUEST ===\n")
		fmt.Printf("[DEBUG] %s %s\n", method, url)
		fmt.Printf("[DEBUG] Headers:\n")
		for key, values := range req.Header {
			for _, value := range values {
				// Mask the API key for security
				if key == "x-api-key" {
					fmt.Printf("[DEBUG] %s: [REDACTED]\n", key)
				} else {
					fmt.Printf("[DEBUG] %s: %s\n", key, value)
				}
			}
		}
		if body != nil {
			// Pretty print the request body for readability
			var prettyBody map[string]interface{}
			if bodyBytes, ok := reqBody.(*bytes.Buffer); ok {
				if err := json.Unmarshal(bodyBytes.Bytes(), &prettyBody); err == nil {
					prettyBytes, err := json.MarshalIndent(prettyBody, "", "  ")
					if err == nil {
						fmt.Printf("[DEBUG] Body:\n%s\n", string(prettyBytes))
					} else {
						fmt.Printf("[DEBUG] Body: %s\n", reqBody)
					}
				} else {
					fmt.Printf("[DEBUG] Body: %s\n", reqBody)
				}
			} else {
				fmt.Printf("[DEBUG] Body: %s\n", reqBody)
			}
		}
		fmt.Printf("[DEBUG] ===================\n")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	if c.config.Debug {
		fmt.Printf("[DEBUG] === HTTP RESPONSE ===\n")
		fmt.Printf("[DEBUG] Status: %s\n", resp.Status)
		fmt.Printf("[DEBUG] Headers:\n")
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("[DEBUG] %s: %s\n", key, value)
			}
		}
		fmt.Printf("[DEBUG] ====================\n")
	}

	return resp, nil
}

// parseJSONResponse parses the JSON response from an HTTP request
func (c *Client) parseJSONResponse(body io.Reader, target interface{}) error {
	responseBody, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	if c.config.Debug {
		fmt.Printf("[DEBUG] Raw Response Body:\n%s\n", string(responseBody))
		fmt.Printf("[DEBUG] ====================\n")
		fmt.Printf("[DEBUG] Pretty-printed Response:\n")

		// Try to pretty print the response for better readability
		var prettyJSON interface{}
		if err := json.Unmarshal(responseBody, &prettyJSON); err == nil {
			prettyBytes, err := json.MarshalIndent(prettyJSON, "", "  ")
			if err == nil {
				fmt.Printf("%s\n", string(prettyBytes))
			} else {
				fmt.Printf("Failed to pretty-print JSON: %v\n", err)
				fmt.Printf("Raw response: %s\n", string(responseBody))
			}
		} else {
			fmt.Printf("Failed to parse JSON for pretty printing: %v\n", err)
			fmt.Printf("Raw response: %s\n", string(responseBody))
		}
		fmt.Printf("[DEBUG] ====================\n")
	}

	if err := json.Unmarshal(responseBody, target); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v (Response body: %s)", err, string(responseBody))
	}

	return nil
}

// ValidateConfig validates the client configuration
func ValidateConfig(config Config) error {
	if config.ServerURL == "" {
		return fmt.Errorf("server URL is required")
	}
	if config.AlbumName == "" {
		return fmt.Errorf("album name is required")
	}
	if config.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	// Ensure server URL doesn't have trailing slash
	config.ServerURL = strings.TrimSuffix(config.ServerURL, "/")

	return nil
}
