import { postgresClientUser } from "../config/DB/DBUser";
import { logger } from "./logger";
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