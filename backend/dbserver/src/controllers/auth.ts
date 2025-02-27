// src/controllers/auth.ts
/**
    Auth Registration Service

    Handles auth registration messages received from the message queue.
    Parses incoming messages, validates data, and inserts the user into the PostgreSQL database.

    Features:
    - Parses message content to extract user details.
    - Executes a database insert query to register the user.
    - Logs success, failure, and error scenarios.

    Dependencies:
    - Logger for structured logging
    - PostgreSQL client (`postgresClientUser`)
    - Query list (`queryList`)

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: 10 Feb 25  
 */


import { logger } from "../utils/logger";
import { postgresClientUser } from "../config/DB/DBUser";
import queryList from "../config/DB/query";

/**
    User Registration Handler

    Processes user registration messages received from the message queue.

    Steps:
    - Parses the incoming message to extract user details.
    - Inserts the user into the PostgreSQL database.
    - Logs the success or failure of the operation.

 */

export const registerUser = async (msg: any): Promise<void> => {
    try {
        const parsedMsg = JSON.parse(msg.content.toString()); // Parse message

        const { id, email, password, type, created_at } = parsedMsg;

        logger.info({ message: "Auth: AddUser query received", id });

        const client = await postgresClientUser.connect();
        const result = await client.query(queryList.addUser, [BigInt(id), email, password, type, created_at, created_at, new Date().toISOString().replace("T", " ").replace("Z", "")]);
        client.release();

        if (result.rowCount === 0) {
            logger.warn({ message: "Auth Register: User insertion failed", id });
            return;
        }

        logger.info({ message: "Auth Register: User successfully registered", id });

    } catch (error: unknown) {
        logger.error({ message: "Error during registration", error: error instanceof Error ? error.message : String(error) });
    }
};
