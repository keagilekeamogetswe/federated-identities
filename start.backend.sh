#!/bin/bash
cd backend

# Default to development if APP_ENV is not set
APP_ENV="${APP_ENV:-development}"

if [ "$APP_ENV" = "development" ]; then
  if [ -f ".secrets.env" ]; then
    echo "Loading local secrets from .secrets.env..."
    set -a
    source .secrets.env
    set +a
  else
    echo "Error: .secrets.env not found!" >&2
    exit 1
  fi
else
  echo "Using GitHub secrets..."
fi

go run .