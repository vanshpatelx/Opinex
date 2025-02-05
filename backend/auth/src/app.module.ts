import { MiddlewareConsumer, Module, NestModule } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { RequestLoggerMiddleware } from './middleware/request-logger.middleware';
import { MetricsModule } from './metrics/metrics.module';
import { LoggerModule } from './logger/logger.module';
import { PubSubModule } from './pubsub/pubsub.module';
import { RedisModule } from './redis/redis.module';
import { DbModule } from './db/db.module';
import { AuthModule } from './auth/auth.module';

@Module({
  imports: [ConfigModule.forRoot(), AuthModule, DbModule, RedisModule, PubSubModule, LoggerModule, MetricsModule],
  providers: [RequestLoggerMiddleware],  // Include middleware if needed
})
export class AppModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(RequestLoggerMiddleware).forRoutes('*');
  }
}

