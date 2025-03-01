#!/bin/bash

# ============================================================
# ✨ Backend Health Check Script
# ============================================================

# 🎨 Define Colors for UI Formatting
GREEN=$(printf "\e[32m")
YELLOW=$(printf "\e[33m")
RED=$(printf "\e[31m")
CYAN=$(printf "\e[36m")
BOLD=$(printf "\e[1m")
RESET=$(printf "\e[0m")

# 🔹 Icons for better visualization
CHECK="✔️"
WARNING="⚠️"
ERROR="❌"
WAITING="⏳"
SUCCESS="🎉"
INFO="🔍"

# ============================================================
# 🏠 Navigate to Project Root
# ============================================================
cd "$(dirname "$0")/.." || { printf "${RED}${ERROR} Failed to change to project root. Please run from within the project.${RESET}\n"; exit 1; }

# 🥒 Health Check URLs
SERVICES=(
    "auth http://localhost:5001/auth/health/"
    "event http://localhost:8000/event/health/"
    "dbserver http://localhost:5003/dbserver/health/"
)

# 🕒 Configuration
MAX_WAIT_TIME=500
WAIT_TIME=0
WAIT_INTERVAL=5

# ============================================================
# 🛡️ Service Health Check Loop
# ============================================================
while true; do
    ALL_READY=true

    for SERVICE in "${SERVICES[@]}"; do
        SERVICE_NAME=$(echo "$SERVICE" | awk '{print $1}')
        SERVICE_URL=$(echo "$SERVICE" | awk '{print $2}')

        # Suppress curl errors and check service health
        if ! curl -sSf "$SERVICE_URL" > /dev/null 2>&1; then
            printf "${YELLOW}${WAITING} Waiting for ${BOLD}$SERVICE_NAME${RESET} at ${BOLD}$SERVICE_URL${RESET}...\n"
            ALL_READY=false
        else
            printf "${GREEN}${CHECK} ${BOLD}$SERVICE_NAME${RESET} is up at ${BOLD}$SERVICE_URL${RESET}!\n"
        fi
    done

    if $ALL_READY; then
        printf "\n${GREEN}${SUCCESS} All backend services are up and running!${RESET}\n"
        exit 0
    fi

    if [[ $WAIT_TIME -ge $MAX_WAIT_TIME ]]; then
        printf "\n${RED}${ERROR} Timed out waiting for backend services after ${MAX_WAIT_TIME} seconds!${RESET}\n"
        exit 1
    fi

    sleep $WAIT_INTERVAL
    WAIT_TIME=$((WAIT_TIME + WAIT_INTERVAL))
done