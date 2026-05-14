#!/bin/bash

# ==============================================================================
# GoNode Setup Script
# Description: Automates the installation of Golang and Node.js for Ubuntu/Debian
# ==============================================================================

# --- Configuration & Colors ---
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BOLD='\033[1m'
NC='\033[0m'

# --- UI Components ---
show_banner() {
    clear
    echo -e "${CYAN}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}${BOLD}          🚀 GoNode ENVIRONMENT SETUP              ${NC}"
    echo -e "${CYAN}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "This script will prepare your system for GoNode.\n"
}

show_footer() {
    echo -e "\n${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}${BOLD}🎉 SETUP COMPLETE!${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "System is now configured for GoNode."
    echo -e "\n${BOLD}To start the application, run:${NC}"
    echo -e "👉 ${CYAN}go run main.go start${NC}"
    echo -e "👉 ${CYAN}./install.sh${NC} (to build the binary)"
    echo -e "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
}

# --- Helper Functions ---
log_step() { echo -e "\n${BLUE}[$1] $2...${NC}"; }
log_info() { echo -e "${YELLOW}💡 $1${NC}"; }
log_error() { echo -e "${RED}❌ $1${NC}"; }
log_success() { echo -e "${GREEN}✅ $1${NC}"; }

# --- Core Logic ---
run_system_update() {
    log_step "1/2" "Updating package list"
    log_info "Ignoring minor errors from third-party repositories..."
    
    # We use || true to prevent the script from exiting if a random repo is broken
    sudo apt update || true
}

run_package_installation() {
    log_step "2/2" "Installing dependencies (Golang, Node.js, Build Tools)"
    log_info "Please enter your password if prompted."
    
    if sudo apt install -y golang nodejs npm build-essential nginx; then
        log_success "All packages installed successfully."
        return 0
    else
        log_error "Failed to install packages. Please check your internet connection."
        return 1
    fi
}

start_setup() {
    show_banner
    
    echo -e "${BOLD}Select your Operating System:${NC}"
    echo -e "1) 🐧 ${BLUE}Ubuntu${NC}"
    echo -e "2) 🌀 ${CYAN}Debian${NC}"
    echo -e "3) ❌ ${RED}Exit${NC}"
    echo -e "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    read -p "Enter choice [1-3]: " choice

    case $choice in
        1|2)
            run_system_update
            if run_package_installation; then
                show_footer
            else
                exit 1
            fi
            ;;
        3)
            echo -e "\n${YELLOW}Setup cancelled.${NC}"
            exit 0
            ;;
        *)
            log_error "Invalid choice. Exiting..."
            exit 1
            ;;
    esac
}

# --- Execute Main ---
start_setup
