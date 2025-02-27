// src/config/DB/DBUser.ts
/**
    PostgreSQL Client (User Database)

    Implements a singleton PostgreSQL client using the `pg` library.

    Features:
    - Establishes a connection pool to the user database.
    - Ensures a single instance of the client (Singleton Pattern).
    - Logs successful connections and errors.
    - Provides a reusable database client for executing queries.

    Dependencies:
    - pg (PostgreSQL client)
    - Logger for structured logging
    - Config module for database credentials

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: February 10, 2025  
 */

import { Pool } from "pg";
import { config } from "../config";
import { logger } from "../../utils/logger";

class PostgresClientUser {
    private static instance: Pool | null = null;

    private constructor() {}

    public static getInstance(): Pool {
        if (!PostgresClientUser.instance) {
            PostgresClientUser.instance = new Pool({
                host: config.dbUser.host,
                port: config.dbUser.port,
                user: config.dbUser.user,
                password: config.dbUser.password,
                database: config.dbUser.database,
            });
            
            PostgresClientUser.instance.on("connect", () => {
                logger.info({
                    message: `Connected to PostgreSQL at ${config.dbUser.host}:${config.dbUser.port}`,
                    service: "postgres",
                });
            });

            PostgresClientUser.instance.on("error", (error) => {
                logger.error({
                    message: "PostgreSQL connection error",
                    service: "postgres",
                    error: error.message,
                });
            });

        }
        return PostgresClientUser.instance;
    }
}

const postgresClientUser = PostgresClientUser.getInstance();
export { postgresClientUser };
