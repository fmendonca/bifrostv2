#!/bin/bash

set -e

API_IMAGE="bifrost-api:0.5"
FRONT_IMAGE="bifrost-frontend:0.6"
API_DIR="./bifrost-api"
FRONT_DIR="./bifrost-frontend"

echo "🔨 Building Bifrost API..."
cd $API_DIR
podman build -t $API_IMAGE .

echo "✅ Bifrost API image built as $API_IMAGE"

echo "🔨 Building Bifrost Frontend..."
cd ../$FRONT_DIR
podman build -t $FRONT_IMAGE .

echo "✅ Bifrost Frontend image built as $FRONT_IMAGE"

echo "🚀 All builds completed successfully!"
