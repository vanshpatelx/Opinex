import * as amqp from 'amqplib';
import { Injectable } from '@nestjs/common';
import { LoggerService } from '../logger/logger.service';

@Injectable()
export class PubSubService {
    private connection: amqp.Connection;
    private channel: amqp.Channel;
    private readonly EXCHANGE_NAME = 'auth_events';
    private readonly RABBITMQ_URL = process.env.RABBITMQ_URL || 'amqp://rabbitmq:5672';

    constructor(private logger: LoggerService) {
        this.init();
    }

    async init() {
        try {
            this.connection = await amqp.connect(this.RABBITMQ_URL);
            this.channel = await this.connection.createChannel();
            await this.channel.assertExchange(this.EXCHANGE_NAME, 'fanout', { durable: true });
            this.logger.info('RabbitMQ Publisher Connected');
        } catch (error) {
            this.logger.error('Failed to connect to RabbitMQ', error);
            // Retry connection after a delay (e.g., 5 seconds)
            setTimeout(() => this.init(), 5000);
        }
    }

    async publishUser(email: string, password: string, userId: bigint) {
        try {
            const message = JSON.stringify({ email, password, userId });
            await this.channel.publish(this.EXCHANGE_NAME, '', Buffer.from(message));
            this.logger.info(`Published to RabbitMQ: ${message}`);
        } catch (error) {
            this.logger.error('Failed to publish message to RabbitMQ', error);
        }
    }

    async consumeUserEvents() {
        try {
            const queue = await this.channel.assertQueue('', { exclusive: true });
            this.channel.bindQueue(queue.queue, this.EXCHANGE_NAME, '');

            this.channel.consume(queue.queue, (msg) => {
                if (msg !== null) {
                    this.logger.info(`Received message from RabbitMQ: ${msg.content.toString()}`);
                    this.channel.ack(msg);
                }
            });

            this.logger.info('RabbitMQ Consumer Listening for Events');
        } catch (error) {
            this.logger.error('Failed to consume RabbitMQ events', error);
        }
    }
}
