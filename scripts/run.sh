#!/bin/bash

# Immich Favorites to Album Runner Script
# This script runs the immich-fave-to-album Docker container with your configured settings

set -e  # Exit on any error

# Default configuration (can be overridden by environment variables)
SERVER_URL="${IMMICH_SERVER:-https://photos.wetsnow.com}"
ALBUM_NAME="${IMMICH_ALBUM:-Favorites}"
IMAGE_TAG="${IMMICH_IMAGE_TAG:-latest}"
DRY_RUN="${IMMICH_DRY_RUN:-false}"
VERBOSE="${IMMICH_VERBOSE:-false}"

# Build the docker run command
DOCKER_CMD="docker run --rm"

# Add environment file if it exists
if [ -f "/home/dant/src/immich-fave-to-album/.env" ]; then
    DOCKER_CMD="$DOCKER_CMD --env-file /home/dant/src/immich-fave-to-album/.env"
else
    echo "Warning: .env file not found. Make sure IMMICH_API_KEY is set in your environment."
    echo "You can create a .env file with: echo 'IMMICH_API_KEY=your-key-here' > .env"
fi

# Add verbose flag if enabled
if [ "$VERBOSE" = "true" ]; then
    DOCKER_CMD="$DOCKER_CMD -e VERBOSE=1"
fi

# Build the final command
DOCKER_CMD="$DOCKER_CMD mrwetsnow/immich-tools:$IMAGE_TAG -server $SERVER_URL -album $ALBUM_NAME"

# Add no-dry-run flag if not in dry run mode
if [ "$DRY_RUN" != "true" ]; then
    DOCKER_CMD="$DOCKER_CMD -no-dry-run"
fi

echo "🚀 Running Immich Favorites to Album..."
echo "📁 Server: $SERVER_URL"
echo "📂 Album: $ALBUM_NAME"
echo "🏷️  Image: mrwetsnow/immich-tools:$IMAGE_TAG"
echo "🔧 Dry Run: $DRY_RUN"
echo "📝 Verbose: $VERBOSE"
echo ""

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "❌ Error: Docker is not installed or not in PATH"
    exit 1
fi

# Execute the command
echo "Executing: $DOCKER_CMD"
echo "----------------------------------------"
exec $DOCKER_CMD
