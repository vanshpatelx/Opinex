import { Injectable, UnauthorizedException } from '@nestjs/common';
import * as bcrypt from 'bcryptjs';
import { JwtService } from '@nestjs/jwt';
import { RedisService } from '../redis/redis.service';
import { DbService } from '../db/db.service';
import { PubSubService } from '../pubsub/pubsub.service';
import { LoggerService } from '../logger/logger.service';
import { trace } from '@opentelemetry/api';
import { generateUniqueId } from '../ID'; // Ensure you have an implementation for generateUniqueId

@Injectable()
export class AuthService {
  private readonly tracer = trace.getTracer('nestjs-auth');

  constructor(
    private readonly jwtService: JwtService,
    private readonly redisService: RedisService,
    private readonly dbService: DbService,
    private readonly pubSubService: PubSubService,
    private readonly logger: LoggerService
  ) { }

  // async register(email: string, password: string) {
  //   return this.tracer.startActiveSpan('auth.register', async (span) => {
  //     try {
  //       const redis = this.redisService.getClient();
  //       const userCacheKey = `user:${email}`;

  //       if (await redis.exists(userCacheKey)) {
  //         this.logger.warn(`User already exists in Redis: ${email}`);
  //         throw new UnauthorizedException('User already exists');
  //       }

  //       const users = await this.dbService.query('SELECT * FROM users WHERE email = $1', [email]);
  //       if (users.length > 0) {
  //         this.logger.warn(`User already exists in DB: ${email}`);
  //         throw new UnauthorizedException('User already exists');
  //       }

  //       const hashedPassword = await bcrypt.hash(password, 10);
  //       const userId = generateUniqueId(); // Ensure this returns a BIGINT-compatible value

  //       // Store in PubSub (assuming it handles BIGINT properly)
  //       await this.pubSubService.publishUser(email, hashedPassword, userId);

  //       // Store in Redis (convert userId to string for compatibility)
  //       await redis.set(
  //         userCacheKey,
  //         JSON.stringify({
  //           email,
  //           password: hashedPassword,
  //           userId: userId.toString(),
  //         }),
  //         { EX: 86400 }
  //       );

  //       // Store in DB
  //       await this.dbService.query(
  //         'INSERT INTO users (id, email, password) VALUES ($1, $2, $3)',
  //         [userId, email, hashedPassword]
  //       );

  //       // JWT does not support BigInt, so convert it to Number
  //       const token = this.jwtService.sign({ email, userId: Number(userId) }, { expiresIn: '7d' });
  //       this.logger.info(`User registered successfully: ${email}`);

  //       return { token };
  //     } catch (error) {
  //       this.logger.error(`Registration failed: ${email}`, { error: error.message });
  //       throw error;
  //     } finally {
  //       span.end();
  //     }
  //   });
  // }

  async register(email: string, password: string) {
    return trace.getTracer('nestjs-auth').startActiveSpan('test-trace', async (span) => {
      console.log('📡 Sending trace to OpenTelemetry');
      span.end();
    });
  }

  async login(email: string, password: string) {
    return this.tracer.startActiveSpan('auth.login', async (span) => {
      try {
        const redis = this.redisService.getClient();
        const userCacheKey = `user:${email}`;
        let user: { email: string; password: string; userId: string } | null;

        const cachedUser = await redis.get(userCacheKey);

        if (cachedUser) {
          user = JSON.parse(cachedUser);
          if (!user.userId) {
            this.logger.error(`Login failed: Missing userId in Redis for ${email}`);
            throw new UnauthorizedException('Invalid credentials');
          }
          user.userId = BigInt(user.userId).toString();
        } else {
          const users = await this.dbService.query('SELECT * FROM users WHERE email = $1', [email]);

          if (!users.length || !users[0].id) {
            this.logger.warn(`Login failed: User not found in DB - ${email}`);
            throw new UnauthorizedException('Invalid credentials');
          }

          user = {
            email: users[0].email,
            password: users[0].password,
            userId: BigInt(users[0].id).toString(),
          };

          // Store user data in Redis for future logins
          await redis.set(userCacheKey, JSON.stringify(user), { EX: 86400 });
        }

        const isPasswordValid = await bcrypt.compare(password, user.password);
        if (!isPasswordValid) {
          this.logger.warn(`Login failed: Incorrect password for ${email}`);
          throw new UnauthorizedException('Invalid credentials');
        }

        const token = this.jwtService.sign({ email, userId: Number(user.userId) }, { expiresIn: '7d' });

        this.logger.info(`User logged in successfully: ${email}`);
        return { token };
      } catch (error) {
        this.logger.error(`Login failed for ${email}: ${error.message}`);
        throw new UnauthorizedException('Login failed');
      } finally {
        span.end();
      }
    });
  }
}
