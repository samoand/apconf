#!/bin/bash

# Function to get the latest Go version
get_latest_go_version() {
  curl -s https://go.dev/VERSION?m=text | head -n 1
}

# Function to download and install Go
install_go() {
  local install_dir=$1
  # local go_version=$(get_latest_go_version)
  local go_version=1.23.0
  local go_tarball="go${go_version}.darwin-arm64.tar.gz"
  local go_url="https://go.dev/dl/${go_tarball}"

  # Check if the directory exists
  if [ ! -d "$install_dir" ]; then
    echo "Creating install directory: $install_dir"
    mkdir -p "$install_dir"
  fi

  # Download the latest Go tarball
  echo "Downloading Go $go_version"
  echo "curl -L -o /tmp/$go_tarball $go_url"
  curl -L -o /tmp/$go_tarball $go_url

  # Extract the tarball to the install directory
  echo "Installing Go to $install_dir"
  tar -C "$install_dir" -xzf /tmp/$go_tarball

  # Clean up the tarball
  rm /tmp/$go_tarball

  echo "Go $go_version installed successfully in $install_dir/go"
}

# Install Go
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_DIR=$(realpath "$SCRIPT_DIR/../external/go1.23.0")
install_go $GO_DIR
