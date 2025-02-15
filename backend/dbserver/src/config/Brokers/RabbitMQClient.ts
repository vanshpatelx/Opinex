import amqp, { Connection, Channel, ConsumeMessage } from "amqplib";
import { config } from "../config";
import { logger } from "../../utils/logger";
import { registerUser } from "../../controllers/auth";

class RabbitMQConsumer {
    private connection: Connection | null = null;
    private channel: Channel | null = null;
    private queue: string;

    constructor(queue: string) {
        this.queue = queue;
    }

    public async connect(): Promise<void> {
        try {
            this.connection = await amqp.connect(config.rabbitmq.url);
            this.channel = await this.connection.createChannel();

            logger.info({
                message: `✅ Connected to RabbitMQ at ${config.rabbitmq.url} for queue: ${this.queue}`,
                service: "rabbitmq",
            });

            this.connection.on("close", () => {
                logger.error({
                    message: `❌ RabbitMQ connection closed for queue: ${this.queue}. Reconnecting...`,
                    service: "rabbitmq",
                });
                this.reconnect();
            });

            this.connection.on("error", (error) => {
                logger.error({
                    message: `⚠️ RabbitMQ connection error for queue: ${this.queue}`,
                    service: "rabbitmq",
                    error: error.message,
                });
            });

            this.consume();
        } catch (error: any) {
            logger.error({
                message: `❌ Failed to connect to RabbitMQ for queue: ${this.queue}`,
                service: "rabbitmq",
                error: error.message,
            });
            this.reconnect();
        }
    }

    private async reconnect(): Promise<void> {
        this.connection = null;
        this.channel = null;
        setTimeout(() => this.connect(), 5000); // Retry after 5 sec
    }

    public async consume(): Promise<void> {
        if (!this.channel) {
            logger.error({
                message: `❌ RabbitMQ channel not initialized for queue: ${this.queue}. Cannot consume.`,
                service: "rabbitmq",
            });
            return;
        }

        this.channel.consume(this.queue, async (msg) => {
            if (msg) {
                try {
                    this.processMessage(msg);
                    this.channel?.ack(msg);
                } catch (error: any) {
                    logger.error({
                        message: `❌ Error processing message for queue: ${this.queue}`,
                        service: "rabbitmq",
                        error: error.message,
                    });
                    this.channel?.nack(msg, false, true); // Requeue message on failure
                }
            }
        });

        logger.info({
            message: `📥 Consuming messages from queue: ${this.queue}`,
            service: "rabbitmq",
        });
    }

    private processMessage(msg: ConsumeMessage): void {
        switch (this.queue) {
            case "auth_registered_queue":
                registerUser(msg);
                console.log("Processing auth_registered_queue message:", msg.content.toString());
                break;
            case "queue_2":
                console.log("Processing queue_2 message:", msg.content.toString());
                break;
            case "queue_3":
                console.log("Processing queue_3 message:", msg.content.toString());
                break;
            default:
                logger.warn({
                    message: `⚠️ No handler found for queue: ${this.queue}`,
                    service: "rabbitmq",
                });
        }
    }
}

export { RabbitMQConsumer };
