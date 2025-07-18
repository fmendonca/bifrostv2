#!/bin/bash

set -e

API_CONTAINER="bifrost-api"
FRONT_CONTAINER="bifrost-frontend"
API_IMAGE="bifrost-api:0.5"
FRONT_IMAGE="bifrost-frontend:0.6"
API_PORT="8080"
FRONT_PORT="3000"
API_URL="http://192.168.86.129:${API_PORT}"
FRONT_SECRET="meuSegredoForte"

# 🛑 Stop and remove old containers if exist
echo "🛑 Stopping and removing old containers if exist..."
podman rm -f $API_CONTAINER || true
podman rm -f $FRONT_CONTAINER || true

# 🚀 Start API container
echo "🚀 Starting Bifrost API on port $API_PORT..."
podman run -d --name $API_CONTAINER --env-file bifrost-api/env -p $API_PORT:8080 $API_IMAGE

# 🚀 Start Frontend container with runtime env vars
echo "🚀 Starting Bifrost Frontend on port $FRONT_PORT..."
podman run -d --name $FRONT_CONTAINER \
  -e REACT_APP_API_URL=$API_URL \
  -e REACT_APP_FRONTEND_SECRET=$FRONT_SECRET \
  -p $FRONT_PORT:8080 \
  $FRONT_IMAGE

# ✅ Status
echo "✅ Bifrost API running at: http://localhost:$API_PORT"
echo "✅ Bifrost Frontend running at: http://localhost:$FRONT_PORT"
