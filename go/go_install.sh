#!/bin/bash

go_version=$1
go_install_dir=$2
os=$(echo "$3" | tr 'A-Z' 'a-z')
arch=$4

# if x86_64 change it to amd64
if [ "$arch" = "x86_64" ]; then
  arch="amd64"
fi

go_tarball="go${go_version}.${os}-${arch}.tar.gz"
go_url="https://go.dev/dl/${go_tarball}"

# Check if the directory exists
if [ ! -d "$go_install_dir" ]; then
    echo "Creating install directory: $go_install_dir"
    mkdir -p "$go_install_dir"
fi

# Download the latest Go tarball
echo "Downloading Go $go_version"
echo "curl -L -o /tmp/$go_tarball $go_url"
curl -L -o /tmp/$go_tarball $go_url

# Extract the tarball to the install directory
echo "Installing Go to $go_install_dir"
tar -C "$go_install_dir" -xzf /tmp/$go_tarball

# Clean up the tarball
rm /tmp/$go_tarball

echo "Go $go_version installed successfully in $go_install_dir/go"
