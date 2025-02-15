import { logger } from "../utils/logger";
import { postgresClientUser } from "../config/DB/DBUser";

const query = {
    addUser: `INSERT INTO users (id, email, password) VALUES ($1, $2, $3);`,
};

export const registerUser = async (msg: any): Promise<void> => {
    try {
        const parsedMsg = JSON.parse(msg.content.toString()); // Parse message
        if (parsedMsg.event !== "UserRegistered") return;

        const { email, password, id } = parsedMsg;

        logger.info({ message: "Auth: AddUser query received", email });

        const client = await postgresClientUser.connect();
        const result = await client.query(query.addUser, [BigInt(id), email, password]);
        client.release();

        if (result.rowCount === 0) {
            logger.warn({ message: "Auth Register: User insertion failed", email });
            return;
        }

        logger.info({ message: "Auth Register: User successfully registered", email });

    } catch (error: unknown) {
        logger.error({ message: "Error during registration", error: error instanceof Error ? error.message : String(error) });
    }
};
