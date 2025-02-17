#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
CERTS_DIR="./certs"
DOMAIN="svc.saltychat.dev"

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Check if running with sudo
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Please run as root (sudo)${NC}"
  exit 1
fi

# Check for required tools
if ! command_exists openssl; then
  echo -e "${RED}openssl is not installed. Please install it:${NC}"
  echo -e "${YELLOW}sudo dnf install openssl${NC}"
  exit 1
fi

# Create certificates directory if it doesn't exist
mkdir -p "$CERTS_DIR"
cd "$CERTS_DIR"

# Generate private key
echo -e "${GREEN}Generating private key...${NC}"
openssl genrsa -out saltychat.key 2048

# Generate certificate
echo -e "${GREEN}Generating certificate...${NC}"
openssl req -new -x509 -sha256 -key saltychat.key -out saltychat.pem -days 365 -subj "/CN=$DOMAIN"

# Set proper permissions
chmod 600 saltychat.key
chmod 644 saltychat.pem

echo -e "${GREEN}Certificates generated in ${CERTS_DIR}!${NC}"
ls -l
