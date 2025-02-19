# Opinex

# Opinex: Opinion Trading System

## Overview

Opinex is a real-time opinion trading platform with a microservices architecture, deployed on Kubernetes for scalability and high performance.

## Architecture

![Opinex Architecture](https://github.com/vanshpatelx/Opinex/blob/main/docs/image/Arch.png)


## Services

- **Auth (TS):** User authentication.
- **Order (Python):** Order processing.
- **Engine (Go):** Trade matching.
- **Event (Python):** Events(Stock) processing.
- **Trade (TS):** Trade execution.
- **DB Server (TS):** Data execution.
- **Holding (Go):** User resource processing.
- **Settlement (TS):** Trade processing.
- **WS Manager (TS):** WebSocket scaling.
- **WS (TS):** Real-time market updates.

## Development & Deployment

- **Docker Compose** for local dev.
- **CI/CD** with unit, integration, and E2E tests.
- **Kubernetes** for cloud deployment.
