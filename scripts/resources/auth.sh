#!/bin/bash
#########################################################################################
#   copy variables from variable.txt, make new env for each services                    #
#   also run local conatiners for dev pupose                                            #
#########################################################################################

source ./../../variables.txt

# Check if PostgreSQL container is already running
if [ "$(docker ps -q -f name=$AUTH_POSTGRES_CONTAINER)" ]; then
    echo "PostgreSQL container is already running."
    echo "Stopping and removing existing PostgreSQL container..."
    docker stop $AUTH_POSTGRES_CONTAINER
    docker rm $AUTH_POSTGRES_CONTAINER
    echo "Starting new PostgreSQL container..."
    docker run -d --name $AUTH_POSTGRES_CONTAINER \
        -e POSTGRES_USER=$AUTH_DB_USER \
        -e POSTGRES_PASSWORD=$AUTH_DB_PASSWORD \
        -e POSTGRES_DB=$AUTH_DB_NAME \
        -p $AUTH_DB_PORT:5432 \
        -v pg_data:/var/lib/postgresql/data \
        --restart unless-stopped \
        $POSTGRES_IMAGE
else
    echo "Starting new PostgreSQL container..."
    docker run -d --name $AUTH_POSTGRES_CONTAINER \
        -e POSTGRES_USER=$AUTH_DB_USER \
        -e POSTGRES_PASSWORD=$AUTH_DB_PASSWORD \
        -e POSTGRES_DB=$AUTH_DB_NAME \
        -p $AUTH_DB_PORT:5432 \
        -v pg_data:/var/lib/postgresql/data \
        --restart unless-stopped \
        $POSTGRES_IMAGE
fi

# Check if Redis container is already running (for Cache)
if [ "$(docker ps -q -f name=$AUTH_REDIS_CONTAINER)" ]; then
    echo "Redis container (Cache) is already running."
    echo "Stopping and removing existing Redis container..."
    docker stop $AUTH_REDIS_CONTAINER
    docker rm $AUTH_REDIS_CONTAINER
    echo "Starting new Redis container (Cache)..."
    docker run -d --name $AUTH_REDIS_CONTAINER \
        -p $AUTH_REDIS_PORT:6379 \
        -e REDIS_PASSWORD=$AUTH_REDIS_PASSWORD \
        -v redis_data:/data \
        --restart unless-stopped \
        $REDIS_IMAGE
else
    echo "Starting Redis container (Cache)..."
    docker run -d --name $AUTH_REDIS_CONTAINER \
        -p $AUTH_REDIS_PORT:6379 \
        -e REDIS_PASSWORD=$AUTH_REDIS_PASSWORD \
        -v redis_data:/data \
        --restart unless-stopped \
        $REDIS_IMAGE
fi

echo "PostgreSQL and Redis (Cache) are running."

# Function to create .env file in a location
create_env_file() {
    local target_path=$1
    local content=$2

    # Go to the project root (assumes the script is inside the project)
    PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || echo "$PWD")"
    FULL_PATH="$PROJECT_ROOT/$target_path"

    # Ensure the target directory exists
    mkdir -p "$(dirname "$FULL_PATH")"

    # Write the .env content
    echo -e "$content" > "$FULL_PATH"

    # Confirm success and print the full path
    if [ -f "$FULL_PATH" ]; then
        echo "✅ .env file successfully created at: $(realpath "$FULL_PATH")"
    else
        echo "❌ Failed to create .env at $FULL_PATH" >&2
        exit 1
    fi
}

# Define the .env content for Service 1 and Service 2 using variables from variables.txt
ENV_CONTENT1="DB_HOST=\"$AUTH_DB_HOST\"\nDB_PORT=\"$AUTH_DB_PORT\"\nDB_USER=\"$AUTH_DB_USER\"\nDB_PASSWORD=\"$AUTH_DB_PASSWORD\"\nDB_NAME=\"$AUTH_DB_NAME\"\nREDIS_HOST=\"$AUTH_REDIS_HOST\"\nREDIS_PORT=\"$AUTH_REDIS_PORT\"\nAUTH_PORT=\"$AUTH_PORT\"\nREDIS_PASSWORD=\"$AUTH_REDIS_PASSWORD\"\nJWT_SECRET=\"$JWT_SECRET\""
ENV_CONTENT2="DB_HOST=\"$SERVICE2_DB_HOST\"\nDB_PORT=\"$SERVICE2_DB_PORT\"\nDB_USER=\"$SERVICE2_DB_USER\"\nDB_PASSWORD=\"$SERVICE2_DB_PASSWORD\"\nDB_NAME=\"$SERVICE2_DB_NAME\"\nREDIS_HOST=\"$SERVICE2_REDIS_HOST\"\nREDIS_PORT=\"$SERVICE2_REDIS_PORT\"\nPORT=\"$SERVICE2_PORT\""

# Define the target locations (relative to the project root)
TARGET1="backend/auth/.env"
TARGET2="backend/service2/.env"

# Create .env files at both locations with their respective content
create_env_file "$TARGET1" "$ENV_CONTENT1"
create_env_file "$TARGET2" "$ENV_CONTENT2"

echo "🚀 All .env files successfully added!"
