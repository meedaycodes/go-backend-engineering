#!/bin/bash

# Deploy a Go binary to a remote Linux server over SSH.
# Usage: update SERVER_IP and KEY_PATH before running.
#
# This script is split into two parts:
#   Part 1 — runs locally on your Mac
#   Part 2 — run manually on the server via SSH

SERVER_IP="<your-server-ip>"
KEY_PATH="~/.ssh/go-backend-dev-key.pem"
REMOTE_USER="ubuntu"

# -----------------------------------------------------------------------------
# PART 1 — Run locally
# -----------------------------------------------------------------------------

# Cross-compile for Linux x86-64. GOOS and GOARCH override the target platform.
# Go's built-in cross-compilation means no special toolchain is needed.
GOOS=linux GOARCH=amd64 go build -o bin/server-linux ./cmd/server

# Verify the binary is a Linux ELF executable, not a macOS Mach-O binary.
file bin/server-linux

# Copy the binary to the server's /tmp directory via SCP (secure copy over SSH).
# /tmp is writable by all users — we move it to its final location on the server.
scp -i $KEY_PATH bin/server-linux $REMOTE_USER@$SERVER_IP:/tmp/server

# Copy the systemd unit file to the server.
scp -i $KEY_PATH goapp.service $REMOTE_USER@$SERVER_IP:/tmp/goapp.service

# -----------------------------------------------------------------------------
# PART 2 — Run on the server via: ssh -i <key> ubuntu@<server-ip>
# -----------------------------------------------------------------------------

# Create a dedicated service account with no home directory and no login shell.
# The service runs as this user — if compromised, the attacker has minimal access.
# sudo useradd --no-create-home --shell /usr/sbin/nologin goapp

# Ensure the standard binary directory exists.
# sudo mkdir -p /usr/local/bin

# Move binary from /tmp to its final location.
# sudo mv /tmp/server /usr/local/bin/server

# Transfer ownership to the service account.
# sudo chown goapp:goapp /usr/local/bin/server

# Set permissions: owner can read and execute, no one else can access it.
# 500 = r-x------
# sudo chmod 500 /usr/local/bin/server

# Install the systemd unit file.
# sudo mv /tmp/goapp.service /etc/systemd/system/goapp.service

# Reload systemd so it picks up the new unit file.
# sudo systemctl daemon-reload

# Enable the service so it starts automatically on boot.
# sudo systemctl enable goapp

# Start the service.
# sudo systemctl start goapp

# Verify it is running.
# sudo systemctl status goapp

# Tail the logs to confirm the service started correctly.
# sudo journalctl -u goapp -n 20
