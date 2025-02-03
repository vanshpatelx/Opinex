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

COPY --from=build-stage /app/backend/service1/dist /app/backend/service1/dist

ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=user
ENV DB_PASSWORD=password
ENV DB_NAME=mydatabase
ENV REDIS_HOST=redis
ENV REDIS_PORT=6379

EXPOSE 3000

CMD ["node", "backend/service1/dist/index.js"]
