// src/config/DB/DBAuth.ts
/**
    PostgreSQL Client (Auth Database)

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

class PostgresClientAuth {
    private static instance: Pool | null = null;

    private constructor() {}

    public static getInstance(): Pool {
        if (!PostgresClientAuth.instance) {
            PostgresClientAuth.instance = new Pool({
                host: config.dbAuth.host,
                port: config.dbAuth.port,
                user: config.dbAuth.user,
                password: config.dbAuth.password,
                database: config.dbAuth.database,
            });
            
            PostgresClientAuth.instance.on("connect", () => {
                logger.info({
                    message: `Connected to PostgreSQL at ${config.dbAuth.host}:${config.dbAuth.port}`,
                    service: "postgres",
                });
            });

            PostgresClientAuth.instance.on("error", (error) => {
                logger.error({
                    message: "PostgreSQL connection error",
                    service: "postgres",
                    error: error.message,
                });
            });

        }
        return PostgresClientAuth.instance;
    }
}

const postgresClientAuth = PostgresClientAuth.getInstance();
export { postgresClientAuth };
