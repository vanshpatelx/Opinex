import { Injectable } from '@nestjs/common';
import * as winston from 'winston';
import { ElasticsearchTransport } from 'winston-elasticsearch';
import { trace, context } from '@opentelemetry/api';
import * as dotenv from 'dotenv';

dotenv.config();

@Injectable()
export class LoggerService {
  private logger: winston.Logger;

  constructor() {
    const esTransport = new ElasticsearchTransport({
      level: 'info',
      clientOpts: {
        node: process.env.ELASTICSEARCH_URL,
      },
      indexPrefix: `${process.env.SERVICE_NAME}-logs`,
    });

    this.logger = winston.createLogger({
      level: 'info',
      format: winston.format.combine(
        winston.format.timestamp(),
        winston.format.json(),
        winston.format.printf(({ level, message, timestamp, meta }) => {
          const span = trace.getSpan(context.active());
          const traceId = span?.spanContext().traceId || 'unknown-trace';

          return JSON.stringify({
            timestamp,
            level,
            service: process.env.SERVICE_NAME,
            traceId,
            message,
            ...(typeof meta === 'object' && meta !== null ? meta : {}),
          });
        })
      ),
      transports: [
        new winston.transports.Console(),
        esTransport,
        new winston.transports.File({ filename: 'logs/error.log', level: 'error' }),
        new winston.transports.File({ filename: 'logs/combined.log' }),
      ],
    });
  }

  log(level: 'info' | 'warn' | 'error' | 'debug', message: string, meta?: Record<string, unknown>) {
    this.logger.log(level, message, meta || {});
  }

  info(message: string, meta?: Record<string, unknown>) {
    this.log('info', message, meta);
  }

  warn(message: string, meta?: Record<string, unknown>) {
    this.log('warn', message, meta);
  }

  error(message: string, meta?: Record<string, unknown>) {
    this.log('error', message, meta);
  }

  debug(message: string, meta?: Record<string, unknown>) {
    this.log('debug', message, meta);
  }
}
