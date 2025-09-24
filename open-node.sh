#!/bin/bash

# This script starts a new (temporary) container using `node:20` image for safe and secure web frontend development
# You don't want to use `node` in your system, do you?

set -e

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker and try again."
    exit 1
fi

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Run the Node.js container with the current directory mounted
docker run --rm -it \
    -v "$SCRIPT_DIR:/app" \
    -w /app \
    -p 5173:5173 \
    node:20 \
    /bin/bash -lc "cd webui && npm ci || npm install && bash"
