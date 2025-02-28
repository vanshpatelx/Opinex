// src/controllers/auth.ts
/**
    Event Registration Service

    Handles event registration messages received from the message queue.
    Parses incoming messages, validates data, and inserts the user into the PostgreSQL database.

    Features:
    - Parses message content to extract event details.
    - Executes a database insert query to register the event.
    - Logs success, failure, and error scenarios.

    Dependencies:
    - Logger for structured logging
    - PostgreSQL client (`postgresClientEvent`)
    - Query list (`queryList`)

    Eventor: Vansh Patel (remotevansh@gmail.com)  
    Date: 28 Feb 25  
 */


import { logger } from "../utils/logger";
import queryList from "../config/DB/query";
import { postgresClientEvent } from "../config/DB/DBEvent";

/**
    Event Registration Handler

    Processes user registration messages received from the message queue.

    Steps:
    - Parses the incoming message to extract user details.
    - Inserts the user into the PostgreSQL database.
    - Logs the success or failure of the operation.

 */

    export const registerEvent = async (msg: any): Promise<void> => {
        try {
            let rawMessage = msg.content.toString().trim();
    
            // ✅ Ensure single quotes are converted to double quotes for valid JSON
            rawMessage = rawMessage.replace(/'/g, '"');
    
            // ✅ Ensure JSON is properly parsed
            const parsedMsg = JSON.parse(rawMessage);
    
            // ✅ Validate and extract properties
            const { id, name, details, status, settlement_time, created_at } = parsedMsg;
    
            if (!id || isNaN(Number(id))) {
                throw new Error(`Invalid ID received: ${id}`);
            }
    
            logger.info({ message: "Event: AddEvent query received", id, settlement_time });
    
            const client = await postgresClientEvent.connect();
            const result = await client.query(queryList.addEvent, [
                BigInt(id), 
                name, 
                details, 
                status, 
                new Date(settlement_time).toISOString().replace("T", " ").replace("Z", ""), 
                new Date(created_at).toISOString().replace("T", " ").replace("Z", ""), 
                new Date(created_at).toISOString().replace("T", " ").replace("Z", ""), 
                new Date().toISOString().replace("T", " ").replace("Z", "")
            ]);
            client.release();
    
            if (result.rowCount === 0) {
                logger.warn({ message: "Event Register: Event insertion failed", id });
                return;
            }
    
            logger.info({ message: "Event Register: Event successfully registered", id });
    
        } catch (error: unknown) {
            logger.error({ 
                message: "Error during registration", 
                error: error instanceof Error ? error.message : String(error),
                rawMessage: msg.content.toString() // Log raw message for debugging
            });
        }
    };
    