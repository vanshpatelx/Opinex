// src/controllers/auth.ts
/**
    Authentication Controller

    Handles user authentication: registration and login.  
    Uses caching (Redis) and pub/sub (RabbitMQ) for efficiency.

    Endpoints:
    - POST /auth/register → Registers a new user.
    - POST /auth/login → Logs in an existing user.

    Dependencies:
    - PostgreSQL, Redis, RabbitMQ
    - Authentication utilities (bcrypt, JWT)
    - Logging and validation

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: February 27, 2025  
 */

import { Request, Response } from "express";
import bcrypt from "bcrypt";
import { z } from 'zod';
import { redisClient } from "../config/Cache/RedisClient";
import { logger } from "../utils/logger";
import { postgresClient } from "../config/DB/db";
import { generateUniqueId } from "../utils/id";
import { generateToken } from "../utils/jwt";
import { rabbitMQClientPromise } from "../config/Brokers/RabbitMQClient";
import queryList from "../config/DB/query";

// Define a schema
const userSchema = z.object({
    email: z.string().email(),
    password: z.string().nonempty(),
});


/**
    User Registration

    Registers a new user by:
    - Validating input.
    - Checking cache and database for existing users.
    - Hashing passwords and storing the user.
    - Sending registration event to RabbitMQ.
    - Returning a JWT on success.

    Returns:
    - 201 Created → User registered successfully.
    - 400 Bad Request → User already exists.
    - 500 Internal Server Error → Unexpected failure.
 */

export const registerUser = async (req: Request, res: Response): Promise<void> => {
    try {
        const { email, password } = userSchema.parse(req.body);
        logger.info({ message: "Auth: Register request received", email });

        const cacheKey = `user:${email}`;
        const existingUserCache = await redisClient.get(cacheKey);

        if (existingUserCache) {
            logger.warn({ message: "Auth Register: User already exists (Cache)", email });
            res.status(400).json({ message: "User already exists" });
            return;
        }

        const client = await postgresClient.connect();
        const result = await client.query(queryList.getUser, [email]);
        client.release();

        logger.info({ message: "Auth Register: Checking DB", email, result });

        if (result.rows.length > 0) {
            logger.warn({ message: "Auth Register: User already exists (DB)", email });

            await redisClient.set(
                cacheKey,
                JSON.stringify({
                    email: result.rows[0].email,
                    password: result.rows[0].password,
                    type: result.rows[0].type,
                    id: result.rows[0].id.toString() // Convert id to string because of BigInt
                }),
                "EX",
                86400
            );

            res.status(400).json({ message: "User already exists" });
            return;
        }

        const userId = generateUniqueId();
        logger.info({ message: "Auth Register: Register approved", email, userId });

        const hashedPassword = await bcrypt.hash(password, 10);

        const userData = { id: userId.toString(), email, password: hashedPassword, type: "USER" };
        await redisClient.set(cacheKey, JSON.stringify(userData), "EX", 86400);
        logger.info({ message: "Auth Register: User cached", email, userId });

        // PubSub
        await (await rabbitMQClientPromise).registerUser(userId, email, hashedPassword, "USER");
        logger.info({ message: "Auth Register: send to pubsub for register", email, userId });

        // Generate JWT token
        const token = generateToken({ id: userId.toString(), email, type: "USER" });
        logger.info({ message: "Auth Register: Token generated", email, userId });

        res.status(201).json({ message: "User registered successfully", token });
        logger.info({ message: "Auth Register: User successfully registered", email, userId });
    } catch (error: unknown) {
        if (error instanceof Error) {
            logger.error({ message: "Error during registration", error: error.message, stack: error.stack });
        } else {
            logger.error({ message: "Error during registration", error: String(error) });
        }
        res.status(500).json({ message: "Internal server error" });
    }
};


interface UserCache {
    email: string,
    password: string,
    id: string,
    type: string
}

/**
    User Login

    Authenticates a user by:
    - Validating input.
    - Checking cache and database for user credentials.
    - Comparing hashed passwords.
    - Returning a JWT on success.

    Returns:
    - 201 Created → User logged in successfully.
    - 401 Unauthorized → Invalid credentials.
    - 500 Internal Server Error → Unexpected failure.
 */

export const loginUser = async (req: Request, res: Response): Promise<void> => {
    try {
        const { email, password } = userSchema.parse(req.body);
        logger.info({ message: "Auth: Login request received", email });

        const cacheKey = `user:${email}`;
        logger.info({ message: "Auth Login: Checking cache", cacheKey });

        const existingUserCache = await redisClient.get(cacheKey);
        let userId: number | null = null;
        let type: string | null = null;

        if (existingUserCache) {
            logger.info({ message: "Auth Login: User found in Cache", email });

            const user: UserCache = JSON.parse(existingUserCache);
            const passwordMatched = await bcrypt.compare(password, user.password);

            if (!passwordMatched) {
                logger.warn({ message: "Auth Login: Invalid credentials (Cache)", email });
                res.status(401).json({ message: "Invalid credentials" });
                return;
            }

            userId = Number(user.id);
            type = user.type;
        } else {
            logger.info({ message: "Auth Login: User not found in Cache, checking DB", email });

            const client = await postgresClient.connect();
            try {
                const result = await client.query(queryList.getUser, [email]);

                if (result.rows.length === 0) {
                    logger.warn({ message: "Auth Login: User not found (DB)", email });
                    res.status(401).json({ message: "Invalid credentials" });
                    return;
                }

                const dbUser = result.rows[0];
                const passwordMatched = await bcrypt.compare(password, dbUser.password);

                if (!passwordMatched) {
                    logger.warn({ message: "Auth Login: Invalid credentials (DB)", email });
                    res.status(401).json({ message: "Invalid credentials" });
                    return;
                }

                // Store in Redis cache
                await redisClient.set(
                    cacheKey,
                    JSON.stringify({
                        email: dbUser.email,
                        password: dbUser.password,
                        id: dbUser.id.toString(), // Convert BigInt to string
                        type: dbUser.type
                    }),
                    "EX",
                    86400
                );

                userId = dbUser.id;
                type = dbUser.type;
                logger.info({ message: "Auth Login: User found in DB and cached", email });
            } finally {
                client.release();
            }
        }

        if (!userId) {
            logger.error({ message: "Auth Login: Unexpected error - userId is null", email });
            res.status(500).json({ message: "Internal server error" });
            return;
        }

        // Generate JWT token
        const token = generateToken({ id: userId.toString(), email, type });
        logger.info({ message: "Auth Login: Token generated", email, userId });

        res.status(200).json({ message: "User login successfully", token });
        logger.info({ message: "Auth: User successfully logged in", email, userId });
    } catch (error: unknown) {
        logger.error({
            message: "Error during login",
            error: error instanceof Error ? error.message : String(error),
            stack: error instanceof Error ? error.stack : undefined
        });
        res.status(500).json({ message: "Internal server error" });
    }
};
