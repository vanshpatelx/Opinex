import { Module } from '@nestjs/common';
import { JwtModule } from '@nestjs/jwt';
import { AuthService } from './auth.service';
import { AuthController } from './auth.controller';
import { RedisService } from '../redis/redis.service';
import { DbService } from '../db/db.service';
import { PubSubService } from '../pubsub/pubsub.service';
import { LoggerService } from '../logger/logger.service';
import { MetricsService } from '../metrics/metrics.service';

@Module({
  imports: [
    JwtModule.register({
      secret: process.env.JWT_SECRET || 'your_secret_key',
      signOptions: { expiresIn: '7d' },
    }),
  ],
  providers: [AuthService, RedisService, DbService, PubSubService, LoggerService, MetricsService],
  controllers: [AuthController],
  exports: [AuthService],
})
export class AuthModule {}
