#!/bin/bash

set -e

# Configuração
API_IMAGE="bifrost-api:0.4"
FRONT_IMAGE="bifrost-frontend:0.6"
API_DIR="./bifrost-api"
FRONT_DIR="./bifrost-frontend"
API_PORT="8080"
FRONT_PORT="3000"
API_URL="http://192.168.86.129:${API_PORT}"

echo "🔨 Building Bifrost API..."
cd $API_DIR
podman build -t $API_IMAGE .

echo "✅ Bifrost API image built as $API_IMAGE"

echo "🔨 Building Bifrost Frontend..."
cd ../$FRONT_DIR
podman build --build-arg REACT_APP_API_URL=$API_URL -t $FRONT_IMAGE .

echo "✅ Bifrost Frontend image built as $FRONT_IMAGE"

echo "🚀 All builds completed successfully!"
