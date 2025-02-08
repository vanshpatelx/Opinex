#!/bin/bash

# Colors for better visibility
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No color

# Check if minimum required arguments are provided
if [ "$#" -lt 3 ]; then
    echo -e "${RED}Usage: $0 <username> <version> <service1> <service2> ...${NC}"
    exit 1
fi

# Extract CLI inputs
USERNAME=$1
VERSION=$2
shift 2  # Remove first two arguments (username and version) from the list

# Iterate over provided service names
for SERVICE in "$@"; do
    OLD_IMAGE="${SERVICE}-service"
    NEW_IMAGE="${USERNAME}/${OLD_IMAGE}:${VERSION}"

    # Check if the image exists
    if docker images | awk '{print $1}' | grep -q "^${OLD_IMAGE}$"; then
        echo -e "${GREEN}Found image: $OLD_IMAGE${NC}"
        
        # Tag the image with the new format
        docker tag "$OLD_IMAGE" "$NEW_IMAGE"
        echo -e "${GREEN}Tagged $OLD_IMAGE as $NEW_IMAGE${NC}"
        
        # Push the image to the registry
        if docker push "$NEW_IMAGE"; then
            echo -e "${GREEN}Successfully pushed $NEW_IMAGE${NC}"
        else
            echo -e "${RED}Failed to push $NEW_IMAGE! Exiting...${NC}"
            exit 1
        fi
    else
        echo -e "${RED}Image $OLD_IMAGE not found, skipping...${NC}"
    fi
done
