#!/bin/bash

# setup.sh - Environment preparation for GoNode

echo "Initializing GoNode Installation..."

# Function to check if a command exists
check_cmd() {
    command -v "$1" >/dev/null 2>&1
}

check_dependencies() {
    local all_installed=true

    if check_cmd go && check_cmd node && check_cmd nginx; then
        echo "Golang, Node.js, and Nginx are already installed! Skipping installation steps..."
        return 0
    fi
    return 1
}

# 1. Update & Basic Tools
sudo apt update
sudo apt install -y curl git build-essential

# 2. Check for Go
if ! check_cmd go; then
    echo "Installing Golang..."
    wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    source ~/.bashrc
else
    echo "Golang already installed."
fi

# 3. Check for Node.js
if ! check_cmd node; then
    echo "Installing Node.js..."
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
    sudo apt install -y nodejs
else
    echo "Node.js already installed."
fi

# 4. Check for Nginx
if ! check_cmd nginx; then
    echo "Installing Nginx..."
    sudo apt install -y nginx
else
    echo "Nginx already installed."
fi

echo "SETUP COMPLETE!"
echo "System is now configured for GoNode."
echo "To start the application, run:"
echo "  go run main.go start"
echo "  ./install.sh (to build the binary)"
