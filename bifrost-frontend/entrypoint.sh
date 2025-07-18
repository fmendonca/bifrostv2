#!/bin/bash

echo "window._env_ = {
  REACT_APP_API_URL: \"${REACT_APP_API_URL}\",
  REACT_APP_FRONTEND_SECRET: \"${REACT_APP_FRONTEND_SECRET}\"
}" > /usr/share/nginx/html/env.js

exec "$@"
