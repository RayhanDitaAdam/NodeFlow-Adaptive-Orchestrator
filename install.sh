#!/bin/bash

# install.sh - Builder for GoNode binary

echo "Initializing GoNode Installation..."

# 1. Build Binary
echo "Building GoNode binary..."
go build -o gonode cmd/gonode/main.go

if [ $? -eq 0 ]; then
    echo "Build Success: 'gonode' binary created."
else
    echo "Build Failed. Please check your Go installation."
    exit 1
fi

# 2. Make Global (Optional)
read -p "Do you want to make 'gonode' a global command? (y/n): " choice
if [[ "$choice" == "y" || "$choice" == "Y" ]]; then
    sudo ln -sf $(pwd)/gonode /usr/local/bin/gonode
    echo "GoNode is now global! You can run 'gonode' from any directory."
fi

echo "Installation finished successfully."
