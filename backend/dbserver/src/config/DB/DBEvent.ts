// src/config/DB/DBEvent.ts
/**
    PostgreSQL Client (Event Database)

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

class PostgresClientEvent {
    private static instance: Pool | null = null;

    private constructor() {}

    public static getInstance(): Pool {
        if (!PostgresClientEvent.instance) {
            PostgresClientEvent.instance = new Pool({
                host: config.dbEvent.host,
                port: config.dbEvent.port,
                user: config.dbEvent.user,
                password: config.dbEvent.password,
                database: config.dbEvent.database,
            });
            
            PostgresClientEvent.instance.on("connect", () => {
                logger.info({
                    message: `Connected to PostgreSQL at ${config.dbEvent.host}:${config.dbEvent.port}`,
                    service: "postgres",
                });
            });

            PostgresClientEvent.instance.on("error", (error) => {
                logger.error({
                    message: "PostgreSQL connection error",
                    service: "postgres",
                    error: error.message,
                });
            });

        }
        return PostgresClientEvent.instance;
    }
}

const postgresClientEvent = PostgresClientEvent.getInstance();
export { postgresClientEvent };
