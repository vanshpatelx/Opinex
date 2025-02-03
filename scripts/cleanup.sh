POSTGRES_CONTAINER="postgres_container"
REDIS_CONTAINER="redis_container"

# Stop and remove the PostgreSQL container
if [ "$(docker ps -q -f name=$POSTGRES_CONTAINER)" ]; then
    echo "Stopping PostgreSQL container..."
    docker stop $POSTGRES_CONTAINER
    docker rm $POSTGRES_CONTAINER
else
    echo "PostgreSQL container is not running."
fi

# Stop and remove the Redis container
if [ "$(docker ps -q -f name=$REDIS_CONTAINER)" ]; then
    echo "Stopping Redis container..."
    docker stop $REDIS_CONTAINER
    docker rm $REDIS_CONTAINER
else
    echo "Redis container is not running."
fi

echo "PostgreSQL and Redis containers have been stopped and removed."
