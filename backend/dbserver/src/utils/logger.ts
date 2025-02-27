// src/utils/logger.ts
/**
    Logger Utility

    Configures and exports a Winston logger for structured logging.

    Features:
    - Logs messages in JSON format with timestamps.
    - Supports logging to both console and a file (`logs/app.log`).

    Dependencies:
    - winston (Logging library)

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: February 10, 2025  
 */
import winston from "winston";

const logFormat = winston.format.combine(
    winston.format.timestamp(),
    winston.format.json()
);

const logger = winston.createLogger({
    level: "info",
    format: logFormat,
    transports: [
        new winston.transports.Console(),
        new winston.transports.File({ filename: "logs/app.log" })
    ],
});

export { logger };
