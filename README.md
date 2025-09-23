# Immich Favorites to Album

A Go script that automatically adds all your favorite photos and videos from Immich to a specified album.

## Features

- Fetches all assets marked as favorite from your Immich server
- Creates a new album with the specified name (if it doesn't exist)
- Adds all favorite assets to the album
- Supports both photos and videos
- Handles existing albums gracefully
- **Dry-run mode enabled by default** - shows what would be done without making actual changes
- **Verbose mode** - detailed logging of all operations
- **Safe by default** - use `-no-dry-run` flag to actually make changes

## Prerequisites

- Go 1.19 or later
- Immich server running and accessible
- API key with appropriate permissions

## Installation

1. Clone or download this repository
2. Navigate to the project directory
3. The script uses only standard library, no additional dependencies required

## Usage

### Set Environment Variable

First, set your Immich API key as an environment variable:

```bash
export IMMICH_API_KEY="your-api-key-here"
```

### Build the Application

```bash
task build
# or
go build -o immich-fave-to-album
```

### Run the Script

```bash
# Using the built binary (dry-run enabled by default)
./immich-fave-to-album -server <immich-server-url> -album <album-name>

# Run without dry-run (actually makes changes)
./immich-fave-to-album -server <immich-server-url> -album <album-name> -no-dry-run

# Run with verbose output
./immich-fave-to-album -server <immich-server-url> -album <album-name> -verbose

# Run with both verbose and no dry-run
./immich-fave-to-album -server <immich-server-url> -album <album-name> -verbose -no-dry-run

# or directly with go run
go run main.go -server <immich-server-url> -album <album-name> -verbose

# or using Task
task run ARGS="-server <immich-server-url> -album <album-name> -verbose"
```

### Examples

```bash
# Local Immich server (dry-run enabled by default)
./immich-fave-to-album -server http://localhost:2283 -album "My Favorites"

# Remote Immich server with verbose output
./immich-fave-to-album -server https://my-immich-server.com -album "Best Photos" -verbose

# With custom port and actually make changes (no dry-run)
./immich-fave-to-album -server http://192.168.1.100:3001 -album "Family Memories" -verbose -no-dry-run

# Using go run with verbose
go run main.go -server http://localhost:2283 -album "My Favorites" -verbose

# Using Task
task run ARGS="-server http://localhost:2283 -album 'My Favorites' -verbose"
```

### Development

```bash
# Run quality checks
task check

# Format code
task fmt

# Run tests (when available)
task test

# Clean build artifacts
task clean

# Build and run checks
task dev

# Show all available tasks
task help
```

### Command Line Options

- `-server`: Immich server URL (required)
- `-album`: Album name to create or use (required)
- `-verbose`: Enable verbose output (optional, default: false)
- `-debug`: Enable debug output with HTTP request/response details (optional, default: false)
- `-no-dry-run`: Disable dry run mode (optional, default: false - dry run is enabled by default)

## How It Works

1. **Authentication**: Uses the `IMMICH_API_KEY` environment variable to authenticate with your Immich server
2. **Configuration**: Parses command line flags including dry-run and verbose options
3. **Asset Discovery**: Searches for all assets marked as favorite
4. **Album Management**: Checks if an album with the specified name already exists
5. **Album Creation**: Creates a new album if one doesn't exist with the given name (dry-run shows what would be created)
6. **Asset Addition**: Adds all favorite assets to the target album (dry-run shows what would be added)

## API Endpoints Used

- `POST /api/search/metadata` - Search for favorite assets with metadata filtering
- `GET /api/albums` - Get existing albums
- `POST /api/albums` - Create new album
- `PUT /api/albums/{id}/assets` - Add assets to album

## Error Handling

The script includes comprehensive error handling for:
- Missing environment variables
- Network connectivity issues
- API authentication failures
- Invalid server responses
- Missing or invalid data

## Notes

- The script will not duplicate assets if they're already in the album
- Album creation includes a timestamp in the description for tracking
- The script uses the "editor" role for album permissions
- Both photos and videos marked as favorites will be included

## Troubleshooting

### Common Issues

1. **"IMMICH_API_KEY environment variable is required"**
   - Make sure you've set the environment variable correctly

2. **"API request failed with status: 401"**
   - Check that your API key is valid and has the necessary permissions

3. **"No favorite assets found"**
   - Verify that you have assets marked as favorites in your Immich library

4. **"failed to get albums"**
   - Ensure your user has permission to view/create albums

### Debug Mode

For more detailed output, you can modify the script to enable debug logging or check the HTTP responses directly.

## Project Structure

```
immich-fave-to-album/
├── main.go              # Main entry point and CLI handling
├── immich/              # Immich API client package
│   ├── client.go        # API client implementation
│   └── types.go         # Data structures and types
├── Taskfile.yml        # Build automation (Task runner)
└── README.md           # This file
```

## Architecture

The application follows a clean architecture pattern:

- **main.go**: Handles command-line interface, configuration validation, and orchestrates the workflow
- **immich package**: Contains all Immich API interaction logic
  - **Client**: Main client struct with HTTP handling and API methods
  - **Types**: All data structures for requests and responses
  - **Validation**: Configuration and input validation

## License

This project is open source and available under the [MIT License](LICENSE).
