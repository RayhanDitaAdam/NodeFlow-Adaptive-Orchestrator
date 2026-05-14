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
    echo -e "${GREEN}✅ GoNode built successfully!${NC}"
    echo -e "🚀 You can now use ${BLUE}./gonode${NC} to run the engine."
    
    # 3. Inform about PATH (Optional)
    echo -e "\n${YELLOW}💡 Pro Tip:${NC}"
    echo -e "To use 'gonode' from anywhere, you can run:"
    echo -e "   ${BLUE}sudo mv gonode /usr/local/bin/${NC}"
else
    echo -e "${RED}❌ Failed to build GoNode. Check main.go for errors.${NC}"
    exit 1
fi
