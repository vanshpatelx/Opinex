// src/config/config.ts
/**
    Configuration Variables

    Centralized storage for environment variables.  
    Loads values from `.env` using `dotenv` for consistency across services.

    Configurations:
    - PostgreSQL → Database credentials.
    - Server → Application port.
    - RabbitMQ → Broker credentials, connection URL, and exchanges.

    Dependencies:
    - dotenv → Loads environment variables from `.env`.

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: February 10, 2025  
 */

    import dotenv from "dotenv";

    dotenv.config();
    
    export const config = {
        port: Number(process.env.PORT) || 5003,
        dbUser: {
            user: process.env.AUTH_DB_USER || "default_user",
            host: process.env.AUTH_DB_HOST || "localhost",
            database: process.env.AUTH_DB_NAME || "auth_db",
            password: process.env.AUTH_DB_PASSWORD || "password",
            port: Number(process.env.AUTH_DB_PORT) || 5432,
        },
        rabbitmq: {
            user: process.env.RABBITMQ_USER || "admin",
            password: process.env.RABBITMQ_PASSWORD || "password",
            host: process.env.RABBITMQ_HOST || "localhost",
            port: Number(process.env.RABBITMQ_PORT) || 5672,
            get url() {
                return `amqp://${this.user}:${this.password}@${this.host}:${this.port}`;
            },
            queues: (process.env.RABBITMQ_QUEUES || "auth_registered_queue").split(","),
        },
    };
    
    export const loadEnv = () => {
        console.log("🚀 Loading Environment Variables\n");
        console.log("🔹 REDIS CONFIGURATION:", {
            REDIS_HOST: process.env.REDIS_HOST || "not set",
            REDIS_PORT: process.env.REDIS_PORT || "not set",
            REDIS_PASSWORD: process.env.REDIS_PASSWORD ? "******" : "not set",
        });
        console.log("🔹 AUTH DATABASE CONFIGURATION:", {
            DB_USER: process.env.AUTH_DB_USER || "not set",
            DB_HOST: process.env.AUTH_DB_HOST || "not set",
            DB_NAME: process.env.AUTH_DB_NAME || "not set",
            DB_PASSWORD: process.env.AUTH_DB_PASSWORD ? "******" : "not set",
            DB_PORT: process.env.AUTH_DB_PORT || "not set",
        });
        
        console.log("🔹 SERVER PORT:", { PORT: process.env.PORT || "5003" });
        console.log("✅ Environment variables loaded!\n");
    };
    