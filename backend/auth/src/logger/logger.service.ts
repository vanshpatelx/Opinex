// import { Injectable } from '@nestjs/common';
// import * as winston from 'winston';
// import { ElasticsearchTransport } from 'winston-elasticsearch';
// import { trace, context } from '@opentelemetry/api';

// @Injectable()
// export class LoggerService {
//   private logger: winston.Logger;

//   constructor() {
//     const esTransport = new ElasticsearchTransport({
//       level: 'info',
//       clientOpts: {
//         node: process.env.ELASTICSEARCH_URL || 'http://localhost:9200',
//       },
//       indexPrefix: 'nestjs-logs',
//     });

//     this.logger = winston.createLogger({
//       level: 'info',
//       format: winston.format.combine(
//         winston.format.timestamp(),
//         winston.format.json(),
//         winston.format.printf(({ level, message, timestamp, meta }) => {
//           const span = trace.getSpan(context.active());
//           const traceId = span?.spanContext().traceId || 'unknown-trace';

//           return JSON.stringify({
//             timestamp,
//             level,
//             message,
//             traceId,
//             ...(typeof meta === 'object' && meta !== null ? meta : {}), // ✅ Ensures meta is an object
//           });          
//         })
//       ),
//       transports: [
//         new winston.transports.Console(),
//         esTransport,
//         new winston.transports.File({ filename: 'logs/error.log', level: 'error' }),
//         new winston.transports.File({ filename: 'logs/combined.log' }),
//       ],
//     });
//   }

//   log(level: 'info' | 'warn' | 'error' | 'debug', message: string, meta?: any) {
//     this.logger.log(level, message, meta);
//   }

//   info(message: string, meta?: any) {
//     this.log('info', message, meta);
//   }

//   warn(message: string, meta?: any) {
//     this.log('warn', message, meta);
//   }

//   error(message: string, meta?: any) {
//     this.log('error', message, meta);
//   }

//   debug(message: string, meta?: any) {
//     this.log('debug', message, meta);
//   }
// }


// import { Injectable } from '@nestjs/common';
// import * as winston from 'winston';
// import { ElasticsearchTransport } from 'winston-elasticsearch';
// import { trace, context } from '@opentelemetry/api';

// @Injectable()
// export class LoggerService {
//   private logger: winston.Logger;

//   constructor() {
//     const esTransport = new ElasticsearchTransport({
//       level: 'info',
//       clientOpts: {
//         node: process.env.ELASTICSEARCH_URL || 'http://localhost:9200',
//       },
//       indexPrefix: 'nestjs-logs',
//     });

//     this.logger = winston.createLogger({
//       level: 'info',
//       format: winston.format.combine(
//         winston.format.timestamp(),
//         winston.format.json(),
//         winston.format.printf(({ level, message, timestamp, meta }) => {
//           const span = trace.getSpan(context.active());
//           const traceId = span?.spanContext().traceId || 'unknown-trace';

//           return JSON.stringify({
//             timestamp,
//             level,
//             message,
//             traceId,
//             ...(typeof meta === 'object' && meta !== null ? meta : {}),
//           });
//         })
//       ),
//       transports: [
//         new winston.transports.Console(),
//         esTransport,
//         new winston.transports.File({ filename: 'logs/error.log', level: 'error' }),
//         new winston.transports.File({ filename: 'logs/combined.log' }),
//       ],
//     });
//   }

//   log(level: 'info' | 'warn' | 'error' | 'debug', message: string, meta?: any) {
//     this.logger.log(level, message, meta);
//   }

//   info(message: string, meta?: any) {
//     this.log('info', message, meta);
//   }

//   warn(message: string, meta?: any) {
//     this.log('warn', message, meta);
//   }

//   error(message: string, meta?: any) {
//     this.log('error', message, meta);
//   }

//   debug(message: string, meta?: any) {
//     this.log('debug', message, meta);
//   }
// }


import { Injectable } from '@nestjs/common';
import * as winston from 'winston';
import { ElasticsearchTransport } from 'winston-elasticsearch';
import { trace, context } from '@opentelemetry/api';

@Injectable()
export class LoggerService {
  private logger: winston.Logger;

  constructor() {
    const esTransport = new ElasticsearchTransport({
      level: 'info',
      clientOpts: {
        node: process.env.ELASTICSEARCH_URL || 'http://elasticsearch:9200',
      },
      indexPrefix: 'nestjs-logs',
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
            message,
            traceId,
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

  log(level: 'info' | 'warn' | 'error' | 'debug', message: string, meta?: any) {
    this.logger.log(level, message, meta);
  }

  info(message: string, meta?: any) {
    this.log('info', message, meta);
  }

  warn(message: string, meta?: any) {
    this.log('warn', message, meta);
  }

  error(message: string, meta?: any) {
    this.log('error', message, meta);
  }

  debug(message: string, meta?: any) {
    this.log('debug', message, meta);
  }
}
