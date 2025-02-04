#!/bin/bash
#########################################################################################
#   check all services running fine or not,                                             #
#   cronjobs to check /health endpoints for each service                                #
#########################################################################################

source ./../variables.txt

MAX_WAIT_TIME=60 # 1 min max time
WAIT_TIME=0
WAIT_INTERVAL=5

# Wait for Backend Services to be Ready (Service 1 and Service 2)
until curl -sSf http://localhost:$AUTH_PORT/health && curl -sSf http://localhost:$AUTH_PORT/health; do
  if [ $WAIT_TIME -ge $MAX_WAIT_TIME ]; then
    echo "❌ Timed out waiting for backend services to be ready after $MAX_WAIT_TIME seconds!"
    exit 1
  fi
  echo "Waiting for backend to be ready on ports $AUTH_PORT and $AUTH_PORT... (Waited $WAIT_TIME seconds)"
  sleep $WAIT_INTERVAL
  WAIT_TIME=$((WAIT_TIME + WAIT_INTERVAL))
done

echo "✅ Both backend services are up and ready!"
