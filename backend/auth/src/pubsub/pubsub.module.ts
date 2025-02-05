import { Module } from '@nestjs/common';
import { LoggerService } from '../logger/logger.service';
import { PubSubService } from './pubsub.service';

@Module({
  providers: [PubSubService, LoggerService],
  exports: [PubSubService],
})
export class PubSubModule {}
