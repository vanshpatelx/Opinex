import { Injectable, NestMiddleware, Logger } from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';
import { MetricsService } from '../metrics/metrics.service';

@Injectable()
export class RequestLoggerMiddleware implements NestMiddleware {
  private readonly logger = new Logger('RequestLogger');

  constructor(private readonly metricsService: MetricsService) {}

  use(req: Request, res: Response, next: NextFunction) {
    const startTime = Date.now();
    res.on('finish', () => {
      const duration = Date.now() - startTime;
      this.metricsService.recordRequest(req.method, req.url, res.statusCode.toString(), duration);
      this.logger.log(`${req.method} ${req.url} - ${res.statusCode} (${duration}ms)`);
    });
    next();
  }
}
