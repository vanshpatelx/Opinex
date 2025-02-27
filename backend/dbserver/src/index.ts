import express, { Express, Request, Response } from "express";
import { config, loadEnv } from "./config/config";
import { initServices } from "./utils/init";
import { logger } from "./utils/logger";
import { RabbitMQConsumer } from "./config/Brokers/RabbitMQConsumer";

loadEnv();

const app: Express = express();
const port = config.port || 5003;

app.use(express.json());
app.use(express.urlencoded({ extended: false }));

app.get('/dbserver/health', (req: Request, res: Response) => {
  res.status(200).json({ success: true, message: "🚀 🚀 Server is running." });
});

async function startConsumers() {
    const queues = ["auth_queue"];

    for (const queue of queues) {
        const consumer = new RabbitMQConsumer(queue);
        await consumer.connect();
    }
}

async function startServer() {
    try {
        await initServices();
        await startConsumers();

        app.listen(port, () => {
            logger.info(`🚀 DBserver is running on port: ${port}`);
        });
    } catch (error) {
        logger.error("❌ Failed to start server:", error);
        process.exit(1);
    }
}

startServer().catch((err) => logger.error("❌ Unexpected error:", err));
