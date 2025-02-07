
#!/bin/bash

# Define colors
GREEN=$(printf "\e[32m")
YELLOW=$(printf "\e[33m")
RED=$(printf "\e[31m")
BLUE=$(printf "\e[34m")
CYAN=$(printf "\e[36m")
BOLD=$(printf "\e[1m")
RESET=$(printf "\e[0m")

# Icons
CHECK="✔️"
WARNING="⚠️"
ERROR="❌"

# Paths
VARIABLES_FILE="../variables.txt"
BACKEND_DIR="../backend"
DOCKER_DIR="../docker"

DOCKER_COMPOSE_FILES=(
    "docker-compose.yaml"
    "docker-compose.resources.yaml"
)

printf "${BOLD}${CYAN}🔍 Detecting services...${RESET}\n"
SERVICES=($(ls -d ${BACKEND_DIR}/*/ 2>/dev/null | xargs -n 1 basename))

if [[ ${#SERVICES[@]} -eq 0 ]]; then
    printf "${RED}${ERROR} No services found in ${BACKEND_DIR}/${RESET}\n"
    exit 1
fi

printf "${GREEN}${CHECK} Found ${#SERVICES[@]} services.${RESET}\n"

for SERVICE in "${SERVICES[@]}"; do
    ENV_FILE="${BACKEND_DIR}/${SERVICE}/.env"
    SERVICE_UPPER=$(echo "$SERVICE" | tr '[:lower:]' '[:upper:]')

    # Check if variables exist for this service
    if ! grep -q "^${SERVICE_UPPER}_" "$VARIABLES_FILE"; then
        printf "${YELLOW}${WARNING} Skipping ${SERVICE}, no matching variables found.${RESET}\n"
        continue
    fi

    # Clear .env file before writing
    printf "${BLUE}🗑️  Clearing $ENV_FILE...${RESET}\n"
    > "$ENV_FILE"

    printf "${GREEN}${CHECK} Generating $ENV_FILE...${RESET}\n"
    while IFS='=' read -r key value; do
        if [[ -n "$key" && "$key" != "#"* ]]; then
            if [[ "$key" == ${SERVICE_UPPER}_* ]]; then
                new_key=${key#${SERVICE_UPPER}_}
                echo "${new_key}=${value}" >> "$ENV_FILE"
            fi
        fi
    done < "$VARIABLES_FILE"

    # Docker Compose Files
    for COMPOSE_FILE in "${DOCKER_COMPOSE_FILES[@]}"; do
        EXAMPLE_FILE="${DOCKER_DIR}/${SERVICE}/example.${COMPOSE_FILE}"
        TARGET_FILE="${DOCKER_DIR}/${SERVICE}/${COMPOSE_FILE}"

        if [[ -f "$EXAMPLE_FILE" ]]; then
            printf "${BLUE}📋 Resetting $TARGET_FILE from $EXAMPLE_FILE...${RESET}\n"
            cp "$EXAMPLE_FILE" "$TARGET_FILE"
        else
            printf "${YELLOW}${WARNING} No example file found for ${COMPOSE_FILE}, skipping reset.${RESET}\n"
            continue
        fi

        if [[ ! -f "$VARIABLES_FILE" ]]; then
            printf "${YELLOW}${WARNING} Variables file not found! Skipping variable replacement.${RESET}\n"
            continue
        fi

        while IFS='=' read -r key value; do
            if [[ -n "$key" && "$key" != "#"* ]]; then
                if [[ "$key" == ${SERVICE_UPPER}_* ]]; then
                    var_name=${key#${SERVICE_UPPER}_}  
                    key_placeholder="\${$var_name}"  

                    if [[ "$OSTYPE" == "darwin"* ]]; then
                        sed -i "" "s#${key_placeholder}#${value}#g" "$TARGET_FILE"
                    else
                        sed -i "s#${key_placeholder}#${value}#g" "$TARGET_FILE"
                    fi
                fi
            fi
        done < "$VARIABLES_FILE"

        printf "${GREEN}${CHECK} Updated $TARGET_FILE.${RESET}\n"
    done
done

printf "${GREEN}${CHECK} Done!${RESET}\n"
