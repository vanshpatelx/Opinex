// src/config/Cache/RedisClient.ts
/**
    Redis Client

    Implements a singleton Redis client using `ioredis`.

    Features:
    - Ensures only one Redis connection is used.
    - Logs successful connection and errors.

    Dependencies:
    - ioredis (Redis client)
    - Logger for structured logging

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: February 2, 2025  
 */


import Redis from "ioredis";
import { config } from "../config";
import { logger } from "../../utils/logger";

class RedisClient {
    private static instance: Redis | null = null;

    private constructor() {}

    public static getInstance(): Redis {
        if (!RedisClient.instance) {
            RedisClient.instance = new Redis({
                host: config.redis.host,
                port: config.redis.port,
                password: config.redis.password,
            });

            RedisClient.instance.on("connect", () => {
                logger.info({
                    message: `Connected to Redis at ${config.redis.host}:${config.redis.port}`,
                    service: "redis",
                });
            });

            RedisClient.instance.on("error", (error) => {
                logger.error({
                    message: "Redis connection error",
                    service: "redis",
                    error: error.message,
                });
            });
        }
        return RedisClient.instance;
    }
}

const redisClient = RedisClient.getInstance();
export { redisClient };
