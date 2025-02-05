// import { Module } from '@nestjs/common';
// import { AuthService } from './auth.service';
// import { AuthController } from './auth.controller';

// @Module({
//   providers: [AuthService],
//   controllers: [AuthController]
// })
// export class AuthModule {}


// import { Module } from '@nestjs/common';
// import { JwtModule } from '@nestjs/jwt';  // ✅ Ensure JwtModule is imported
// import { AuthService } from './auth.service';
// import { AuthController } from './auth.controller';
// import { RedisService } from '../redis/redis.service';
// import { DbService } from '../db/db.service';
// import { PubSubService } from '../pubsub/pubsub.service';
// import { LoggerService } from '../logger/logger.service';

// @Module({
//   imports: [
//     JwtModule.register({
//       secret: process.env.JWT_SECRET || 'your_secret_key', // ✅ Ensure secret is set
//       signOptions: { expiresIn: '7d' },  // ✅ Set token expiration
//     }),
//   ],
//   providers: [AuthService, RedisService, DbService, PubSubService, LoggerService],
//   controllers: [AuthController],
//   exports: [AuthService],  // ✅ Export AuthService if needed in other modules
// })
// export class AuthModule {}


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
