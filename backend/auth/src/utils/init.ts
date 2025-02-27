// src/utils/init.ts
/**
    Service Health Checker

    This module checks the readiness of essential services:
    - PostgreSQL Database  
    - Redis Cache  
    - RabbitMQ Message Broker  

    The system retries failed connections up to 5 times with a delay before exiting.  

    Author: Vansh Patel (remotevansh@gmail.com)  
    Last Updated Date: February 2, 2025  
 */  


import { postgresClient } from "../config/DB/db";
import { redisClient } from "../config/Cache/RedisClient";
import { logger } from "./logger";
import amqp from "amqplib";
import { config } from "../config/config";


/**
    PostgreSQL Readiness Check

    Executes a simple query to verify if the PostgreSQL database is available.  
    Logs an error if the database is not accessible.

    Returns:
    - Promise<boolean> → Resolves to true if PostgreSQL is ready, otherwise false.  
 */  

async function checkPostgresConnection(): Promise<boolean> {
    return new Promise((resolve) => {
        postgresClient.query("SELECT 1", (err) => {
            if (err) {
                logger.error("PostgreSQL is not ready:", err.message);
                resolve(false);
            } else {
                logger.info("PostgreSQL is ready.");
                resolve(true);
            }
        });
    });
}


/**
    Redis Readiness Check

    Sends a PING command to verify if the Redis cache is available.  
    Logs an error if Redis is not accessible.

    Returns:
    - Promise<boolean> → Resolves to true if Redis is ready, otherwise false.  
 */  

async function checkRedisConnection(): Promise<boolean> {
    return new Promise((resolve) => {
        redisClient.ping((err, result) => {
            if (err || result !== "PONG") {
                logger.error("Redis is not ready:", err?.message);
                resolve(false);
            } else {
                logger.info("Redis is ready.");
                resolve(true);
            }
        });
    });
}


/**
    RabbitMQ Readiness Check

    Establishes a temporary connection to verify if RabbitMQ is available.  
    Logs an error if RabbitMQ is not accessible.

    Returns:
    - Promise<boolean> → Resolves to true if RabbitMQ is ready, otherwise false.  
 */  

async function checkRabbitMQConnection(): Promise<boolean> {
    try {
        const connection = await amqp.connect(config.rabbitmq.url);
        await connection.close();
        logger.info("RabbitMQ is ready.");
        return true;
    } catch (error: any) {
        logger.error("RabbitMQ is not ready:", error.message);
        return false;
    }
}


/**
    Service Initialization

    Initializes required services: PostgreSQL, Redis, and RabbitMQ.  
    Retries up to 5 times with a 5-second delay if any service is not ready.  
    Logs success if all services are ready, otherwise throws an error.  
 */  

export async function initServices() {
    let retries = 5;
    while (retries > 0) {
        const dbReady = await checkPostgresConnection();
        const cacheReady = await checkRedisConnection();
        const rabbitMQReady = await checkRabbitMQConnection();

        if (dbReady && cacheReady && rabbitMQReady) {
            logger.info("✅ All services are ready!");
            return;
        }

        logger.warn(`Retrying service initialization (${retries} attempts left)...`);
        await new Promise((res) => setTimeout(res, 5000)); // Wait before retrying
        retries--;
    }

    throw new Error("❌ Services failed to initialize. Exiting...");
}