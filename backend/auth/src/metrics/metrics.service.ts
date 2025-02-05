import { Injectable, OnModuleInit } from '@nestjs/common';
import { register, Counter, Histogram } from 'prom-client';

@Injectable()
export class MetricsService implements OnModuleInit {
  private requestCount: Counter<string>;
  private requestDuration: Histogram<string>;

  onModuleInit() {
    // Set global labels unconditionally
    register.setDefaultLabels({
      service: process.env.SERVICE_NAME,
    });

    // Check if the metric 'http_requests_total' already exists
    const existingRequestCount = register.getSingleMetric('http_requests_total') as Counter<string>;
    if (existingRequestCount) {
      this.requestCount = existingRequestCount;
    } else {
      this.requestCount = new Counter({
        name: 'http_requests_total',
        help: 'Total number of HTTP requests',
        labelNames: ['method', 'route', 'status'],
      });
    }

    // Check if the metric 'http_request_duration_seconds' already exists
    const existingRequestDuration = register.getSingleMetric('http_request_duration_seconds') as Histogram<string>;
    if (existingRequestDuration) {
      this.requestDuration = existingRequestDuration;
    } else {
      this.requestDuration = new Histogram({
        name: 'http_request_duration_seconds',
        help: 'Duration of HTTP requests in seconds',
        labelNames: ['method', 'route', 'status'],
        buckets: [0.1, 0.5, 1, 3, 5],
      });
    }
  }

  recordRequest(method: string, route: string, status: string, duration: number) {
    this.requestCount.labels(method, route, status).inc();
    this.requestDuration.labels(method, route, status).observe(duration);
  }

  async getMetrics(): Promise<string> {
    return register.metrics();
  }
}
