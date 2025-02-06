import { Pool } from "pg";
import { config } from "../config";
import { logger } from "../../utils/logger";

const pool = new Pool({
    host: config.db.host,
    port: config.db.port,
    user: config.db.user,
    password: config.db.password,
    database: config.db.database
});

pool.on("connect", () => {
    logger.info({ message: "Connected to PostgreSQL database" });
});

pool.on("error", (err) => {
    logger.error({ message: "PostgreSQL error", error: err.message });
});

export default pool;
