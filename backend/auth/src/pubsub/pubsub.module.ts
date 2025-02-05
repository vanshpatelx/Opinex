import { Module } from '@nestjs/common';
import { LoggerService } from '../logger/logger.service';  // ✅ Import LoggerService
import { PubSubService } from './pubsub.service';  // Assuming PubSubService exists here

@Module({
  imports: [], // Add any modules that provide dependencies here
  providers: [PubSubService, LoggerService], // Ensure LoggerService is part of providers
  exports: [PubSubService], // Export if needed elsewhere
})
export class PubSubModule {}
