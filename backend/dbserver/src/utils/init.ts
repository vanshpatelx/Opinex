// src/utils/init.ts
/**
    Service Initialization

    Ensures PostgreSQL is ready before proceeding with service startup.
    Retries connection checks up to 5 times before failing.

    Features:
    - `checkPostgresConnection()`: Verifies database connectivity.
    - `initServices()`: Attempts to initialize services with retry logic.

    Dependencies:
    - PostgreSQL client (`postgresClientUser`)
    - Logger for structured logging

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: 10 Feb 25  
 */

    
import { postgresClientUser } from "../config/DB/DBUser";
import { logger } from "./logger";


/**
    PostgreSQL Readiness Check

    Executes a simple query to verify if the PostgreSQL database is available.  
    Logs an error if the database is not accessible.

    Returns:
    - Promise<boolean> → Resolves to true if PostgreSQL is ready, otherwise false.  
 */  
async function checkPostgresConnection(): Promise<boolean> {
    return new Promise((resolve) => {
        postgresClientUser.query("SELECT 1", (err) => {
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
    Service Initialization

    Initializes required services: PostgreSQL.  
    Retries up to 5 times with a 5-second delay if any service is not ready.  
    Logs success if all services are ready, otherwise throws an error.  
 */  
export async function initServices() {
    let retries = 5;
    while (retries > 0) {
        const dbReady = await checkPostgresConnection();

        if (dbReady) {
            logger.info("✅ All services are ready!");
            return;
        }

        logger.warn(`Retrying service initialization (${retries} attempts left)...`);
        await new Promise((res) => setTimeout(res, 5000)); // Wait before retrying
        retries--;
    }

    throw new Error("❌ Services failed to initialize. Exiting...");
}