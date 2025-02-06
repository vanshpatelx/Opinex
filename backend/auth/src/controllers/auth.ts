import { Request, Response } from "express";
import bcrypt from "bcrypt";
import { z } from 'zod';
import { redisClient } from "../config/Cache/RedisClient";
import { logger } from "../utils/logger";
import pool from "../config/DB/db";
import { generateUniqueId } from "../utils/id";
import { generateToken } from "../utils/jwt";

const users: { email: string; password: string }[] = [];

// Define a schema
const userSchema = z.object({
    email: z.string().email(),
    password: z.string()
});


/**
 * @description Handles user registration  
 * @route POST /auth/register  
 * 
 * ✔ Validates request payload using Zod  
 * ✔ Checks if the user already exists in the cache  
 * ✔ Queries the database for an existing user  
 * ✔ Generates a unique user ID  
 * ✔ Stores user data in the cache  
 * ✔ Publishes user registration event to Pub/Sub  
 */


export const registerUser = async (req: Request, res: Response): Promise<void> => {
    try {
        const { email, password } = userSchema.parse(req.body);
        logger.info({ message: "Auth: Register request received", email });

        const cacheKey = `user:${email}`;
        const existingUserCache = await redisClient.get(cacheKey);

        if (existingUserCache) {
            logger.warn({ message: "Auth: User already exists (Cache)", email });
            res.status(400).json({ message: "User already exists" });
            return;
        }

        const client = await pool.connect();
        const result = await client.query("SELECT * FROM users WHERE email = $1", [email]);
        client.release();

        if (result.rows.length > 0) {
            logger.warn({ message: "Auth: User already exists (DB)", email });

            await redisClient.set(
                cacheKey,
                JSON.stringify({
                    email: result.rows[0].email,
                    password: result.rows[0].password,
                    id: result.rows[0].id.toString() // Convert id to string because of BigInt
                }),
                "EX",
                86400
            );

            res.status(400).json({ message: "User already exists" });
            return;
        }

        const userId = generateUniqueId();
        logger.info({ message: "Auth: Register approved", email, userId });

        const hashedPassword = await bcrypt.hash(password, 10);

        const userData = { id: userId, email, password: hashedPassword };
        await redisClient.set(cacheKey, JSON.stringify(userData), "EX", 86400);
        logger.info({ message: "Auth: User cached", email, userId });

        // PubSub logic (if applicable)
        logger.info({ message: "Auth: send to pubsub for register", email, userId });

        // Generate JWT token
        const token = generateToken({ id: userId, email });
        logger.info({ message: "Auth: Token generated", email, userId });

        res.status(201).json({ message: "User registered successfully", token });
        logger.info({ message: "Auth: User succesftly registered", email, userId });
    } catch (error) {
        logger.error({ message: "Error during registration", error: error.message });
        res.status(500).json({ message: "Internal server error" });
    }
};



/**
 * @desc Login a user
 * @route POST /auth/login
 */
export const loginUser = async (req: Request, res: Response): Promise<void> => {
    const { email, password } = req.body;

    // Find user
    const user = users.find(user => user.email === email);
    if (!user) {
        res.status(401).json({ message: "Invalid email or password" });
        return;
    }

    // Compare hashed password
    const isMatch = await bcrypt.compare(password, user.password);
    if (!isMatch) {
        res.status(401).json({ message: "Invalid email or password" });
        return;
    }

    res.status(200).json({ message: "Login successful" });
};
