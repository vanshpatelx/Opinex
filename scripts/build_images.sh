# #!/bin/bash

# # Move to the script's parent directory (assumes `scripts/` is directly under the root)
# cd "$(dirname "$0")/.." || { echo "Failed to change to project root"; exit 1; }

# # Define colors
# GREEN=$(printf "\e[32m")
# YELLOW=$(printf "\e[33m")
# RED=$(printf "\e[31m")
# CYAN=$(printf "\e[36m")
# BOLD=$(printf "\e[1m")
# RESET=$(printf "\e[0m")

# # Icons
# CHECK="✔️"
# WARNING="⚠️"
# ERROR="❌"
# BUILD="🛠️"

# # Ensure script is run from the root directory
# if [[ ! -d "backend" || ! -d "docker" ]]; then
#     printf "${RED}${ERROR} Script must be run from the project root!${RESET}\n"
#     exit 1
# fi

# # Version (default: latest if not provided)
# VERSION=${1:-latest}

# printf "${BOLD}${CYAN}🔍 Detecting services in backend/...${RESET}\n"
# SERVICES=($(ls -d backend/*/ 2>/dev/null | xargs -n 1 basename))

# if [[ ${#SERVICES[@]} -eq 0 ]]; then
#     printf "${RED}${ERROR} No services found in backend/${RESET}\n"
#     exit 1
# fi

# printf "${GREEN}${CHECK} Found ${#SERVICES[@]} services.${RESET}\n"

# # Loop through each service and build its image
# for SERVICE in "${SERVICES[@]}"; do
#     SERVICE_DOCKERFILE="docker/${SERVICE}/Dockerfile"

#     if [[ ! -f "$SERVICE_DOCKERFILE" ]]; then
#         printf "${YELLOW}${WARNING} No Dockerfile found for ${SERVICE} in docker/${SERVICE}, skipping.${RESET}\n"
#         continue
#     fi

#     IMAGE_NAME="${SERVICE}-service:${VERSION}"
    
#     printf "${CYAN}${BUILD} Building image for ${SERVICE}...${RESET}\n"
#     docker build -f "$SERVICE_DOCKERFILE" -t "$IMAGE_NAME" .

#     if [[ $? -eq 0 ]]; then
#         printf "${GREEN}${CHECK} Successfully built ${IMAGE_NAME}${RESET}\n"
#     else
#         printf "${RED}${ERROR} Failed to build ${SERVICE}${RESET}\n"
#     fi
# done

# printf "${GREEN}${CHECK} Done! All images built.${RESET}\n"


#!/bin/bash

# Move to the script's parent directory (assumes `scripts/` is directly under the root)
cd "$(dirname "$0")/.." || { echo "Failed to change to project root"; exit 1; }

# Define colors
GREEN=$(printf "\e[32m")
YELLOW=$(printf "\e[33m")
RED=$(printf "\e[31m")
CYAN=$(printf "\e[36m")
BOLD=$(printf "\e[1m")
RESET=$(printf "\e[0m")

# Icons
CHECK="✔️"
WARNING="⚠️"
ERROR="❌"
BUILD="🛠️"

# Ensure script is run from the root directory
if [[ ! -d "backend" || ! -d "docker" ]]; then
    printf "${RED}${ERROR} Script must be run from the project root!${RESET}\n"
    exit 1
fi

# Version (default: latest if not provided)
VERSION=${1:-latest}

printf "${BOLD}${CYAN}🔍 Detecting services in backend/...${RESET}\n"
SERVICES=($(ls -d backend/*/ 2>/dev/null | xargs -n 1 basename))

if [[ ${#SERVICES[@]} -eq 0 ]]; then
    printf "${RED}${ERROR} No services found in backend/${RESET}\n"
    exit 1
fi

printf "${GREEN}${CHECK} Found ${#SERVICES[@]} services.${RESET}\n"

VARIABLES_FILE="variables.txt"

# Detect OS and set proper sed syntax
if [[ "$OSTYPE" == "darwin"* ]]; then
    SED_CMD="sed -i ''"
else
    SED_CMD="sed -i"
fi

# Loop through each service and build its image
for SERVICE in "${SERVICES[@]}"; do
    SERVICE_DOCKERFILE="docker/${SERVICE}/Dockerfile"

    if [[ ! -f "$SERVICE_DOCKERFILE" ]]; then
        printf "${YELLOW}${WARNING} No Dockerfile found for ${SERVICE} in docker/${SERVICE}, skipping.${RESET}\n"
        continue
    fi

    IMAGE_NAME="${SERVICE}-service:${VERSION}"
    
    printf "${CYAN}${BUILD} Building image for ${SERVICE}...${RESET}\n"
    docker build -f "$SERVICE_DOCKERFILE" -t "$IMAGE_NAME" .

    if [[ $? -eq 0 ]]; then
        printf "${GREEN}${CHECK} Successfully built ${IMAGE_NAME}${RESET}\n"

        # Convert service name to uppercase with underscore prefix (e.g., auth -> AUTH_IMAGE)
        IMAGE_VAR="$(echo "${SERVICE}_IMAGE" | tr '[:lower:]' '[:upper:]')"

        # Update variables.txt dynamically
        if [[ -f "$VARIABLES_FILE" ]]; then
            $SED_CMD "s|^${IMAGE_VAR}=.*|${IMAGE_VAR}=\"$IMAGE_NAME\"|" "$VARIABLES_FILE"
            printf "${GREEN}${CHECK} Updated ${IMAGE_VAR} in ${VARIABLES_FILE} to ${IMAGE_NAME}${RESET}\n"
        else
            printf "${RED}${ERROR} ${VARIABLES_FILE} not found! Skipping update.${RESET}\n"
        fi
    else
        printf "${RED}${ERROR} Failed to build ${SERVICE}${RESET}\n"
    fi
done

printf "${GREEN}${CHECK} Done! All images built and updated in ${VARIABLES_FILE}.${RESET}\n"
