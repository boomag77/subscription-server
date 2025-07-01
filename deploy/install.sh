#!/bin/bash

set -e

APP_DIR="/opt/subscription-server"
BIN_NAME="server"

echo "Cloning or updating source..."
if [ ! -d "$APP_DIR/.git" ]; then
  git clone https://github.com/boomag77/subscription-server.git "$APP_DIR"
else
  cd "$APP_DIR"
  git reset --hard
  git pull
fi

echo "🔨 Building binary..."
cd "$APP_DIR"
go build -o $BIN_NAME ./cmd/server

echo "Restarting service..."
sudo systemctl restart subscription-server.service

echo "Deploy complete"
