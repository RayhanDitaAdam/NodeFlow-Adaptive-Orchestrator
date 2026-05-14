#!/bin/bash

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🛠️  Initializing GoNode Installation...${NC}"

# 1. Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Error: Go is not installed. Please install Go first.${NC}"
    exit 1
fi

# 2. Check if Node.js is installed (required by app.js)
if ! command -v node &> /dev/null; then
    echo -e "${YELLOW}⚠️  Warning: Node.js is not found. GoNode requires Node.js to run app.js.${NC}"
fi

# 3. Build the binary
echo -e "${YELLOW}📦 Building GoNode binary...${NC}"
go build -o gonode cmd/gonode/main.go

if [ $? -eq 0 ]; then
    chmod +x gonode
    echo -e "${GREEN}✅ Build successful!${NC}"
    
    # 4. Make it global (Optional)
    echo -e "\n${YELLOW}🌍 Do you want to make 'gonode' command global? (y/n)${NC}"
    read -r make_global
    if [[ "$make_global" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        echo -e "${YELLOW}🔑 Requesting sudo permission to create symlink...${NC}"
        sudo ln -sf "$(pwd)/gonode" /usr/local/bin/gonode
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}🚀 Success! You can now run 'gonode' from any directory.${NC}"
        else
            echo -e "${RED}❌ Failed to create symlink.${NC}"
        fi
    fi
    
    echo -e "\n${GREEN}Usage: ./gonode start or simply 'gonode start' if global.${NC}"
else
    echo -e "${RED}❌ Failed to build GoNode. Check main.go for errors.${NC}"
    exit 1
fi
