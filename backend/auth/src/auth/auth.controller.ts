import { Controller, Post, Body, Get, Req, BadRequestException, UnauthorizedException } from '@nestjs/common';
import { Request } from 'express';
import { z } from 'zod';
import { AuthService } from './auth.service';
import { LoggerService } from '../logger/logger.service';
import { MetricsService } from '../metrics/metrics.service';

const RegisterSchema = z.object({
  email: z.string().email('Invalid email format').nonempty('Email is required'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
});

@Controller('auth')
export class AuthController {
  constructor(
    private readonly authService: AuthService,
    private readonly logger: LoggerService,
    private readonly metricsService: MetricsService
  ) {}

  @Post('register')
  async register(
    @Body() body: { email: string; password: string },
    @Req() req: Request
  ) {
    const userIp = req.ip;
    try {
      this.logger.info(`Registration attempt from IP: ${userIp}, Email: ${body.email}`);
      RegisterSchema.parse(body);
      this.logger.info(`Validation successful for Email: ${body.email}`);

      const result = await this.authService.register(body.email, body.password);
      this.logger.info(`Registration successful for Email: ${body.email}`);
      return result;
    } catch (error) {
      this.handleError(error, 'Registration', body.email, userIp);
    }
  }

  @Post('login')
  async login(
    @Body() body: { email: string; password: string },
    @Req() req: Request
  ) {
    const userIp = req.ip;
    try {
      this.logger.info(`Login attempt from IP: ${userIp}, Email: ${body.email}`);
      RegisterSchema.parse(body); // Same validation as registration
      this.logger.info(`Validation successful for Email: ${body.email}`);

      const result = await this.authService.login(body.email, body.password);
      this.logger.info(`Login successful for Email: ${body.email}`);
      return result;
    } catch (error) {
      this.handleError(error, 'Login', body.email, userIp);
    }
  }

  @Get('/metrics')
  getMetrics() {
    return this.metricsService.getMetrics();
  }

  private handleError(error: any, action: string, email: string, userIp: string) {
    if (error instanceof z.ZodError) {
      this.logger.error(`${action} validation failed from IP: ${userIp}, Email: ${email}`, { errors: error.errors });
      throw new BadRequestException(error.errors);
    }

    this.logger.error(`${action} failed from IP: ${userIp}, Email: ${email}`, { error: error.message });
    throw new UnauthorizedException(error.message);
  }
}
