# Stage 1: Build stage
FROM node:18 AS build-stage

WORKDIR /app

# Copy all files from the root directory to the /app directory in the container
COPY ../.. .

RUN npm install

RUN npm install -g turbo

RUN turbo build

# Stage 2: Production image
FROM node:18-slim AS production-stage

WORKDIR /app

# Copy the necessary files from the build stage
COPY --from=build-stage /app /app

COPY --from=build-stage /app/backend/dbserver/dist /app/backend/dbserver/dist

EXPOSE 5003

CMD ["node", "backend/dbserver/dist/index.js"]
