// src/config/DB/db.ts
/**
    PostgreSQL Database Client

    Implements a singleton PostgreSQL client using `pg.Pool` for efficient database connections.

    Features:
    - Singleton pattern ensures only one connection pool is created.
    - Logs successful connection and errors.

    Dependencies:
    - pg (PostgreSQL client)
    - Logger for structured logging

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: February 2, 2025  
 */


import { Pool } from "pg";
import { config } from "../config";
import { logger } from "../../utils/logger";

class PostgresClient {
    private static instance: Pool | null = null;

    private constructor() {}

    public static getInstance(): Pool {
        if (!PostgresClient.instance) {
            PostgresClient.instance = new Pool({
                host: config.db.host,
                port: config.db.port,
                user: config.db.user,
                password: config.db.password,
                database: config.db.database,
            });
            
            PostgresClient.instance.on("connect", () => {
                logger.info({
                    message: `Connected to PostgreSQL at ${config.db.host}:${config.db.port}`,
                    service: "postgres",
                });
            });

            PostgresClient.instance.on("error", (error) => {
                logger.error({
                    message: "PostgreSQL connection error",
                    service: "postgres",
                    error: error.message,
                });
            });

        }
        return PostgresClient.instance;
    }
}

const postgresClient = PostgresClient.getInstance();
export { postgresClient };
