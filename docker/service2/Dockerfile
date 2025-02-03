# Stage 1: Build Go service
FROM golang:1.23 AS go-build

WORKDIR /app/backend/service2

COPY ../.. /app

RUN go mod tidy

RUN go build -o bin/main ./main.go

# Stage 2: Production image
FROM golang:1.23 AS production-stage

WORKDIR /app

# Copy compiled Go binary from the build stage
COPY --from=go-build /app/backend/service2/bin /app/backend/service2/bin

# Copy the .env file from the build stage
COPY --from=go-build /app/backend/service2/.env .env

ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=user
ENV DB_PASSWORD=password
ENV DB_NAME=mydatabase
ENV REDIS_HOST=redis
ENV REDIS_PORT=6379

EXPOSE 8080

CMD ["/app/backend/service2/bin/main"]
