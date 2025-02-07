#!/bin/bash
#########################################################################################
#   Cleaup of all resources used for dev purpose                                        #
#########################################################################################

source ./../variables.txt

if [ "$(docker ps -q -f name=$AUTH_POSTGRES_CONTAINER)" ]; then
    echo "Stopping PostgreSQL container..."
    docker stop $POSTGRES_CONTAINER
    docker rm $POSTGRES_CONTAINER
else
    echo "PostgreSQL container is not running."
fi

if [ "$(docker ps -q -f name=$AUTH_POSTGRES_CONTAINER)" ]; then
    echo "Stopping Redis container..."
    docker stop $REDIS_CONTAINER
    docker rm $REDIS_CONTAINER
else
    echo "Redis container is not running."
fi

echo "PostgreSQL and Redis containers have been stopped and removed."
