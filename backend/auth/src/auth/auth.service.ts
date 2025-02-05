// import { Injectable, UnauthorizedException } from '@nestjs/common';
// import * as bcrypt from 'bcryptjs';
// import { JwtService } from '@nestjs/jwt';
// import { RedisService } from '../redis/redis.service';
// import { DbService } from '../db/db.service';
// import { PubSubService } from '../pubsub/pubsub.service';
// import { LoggerService } from '../logger/logger.service';
// import { trace } from '@opentelemetry/api';
// import { generateUniqueId } from '../ID';

// @Injectable()
// export class AuthService {
//     private readonly tracer = trace.getTracer('nestjs-auth');

//     constructor(
//         private readonly jwtService: JwtService,
//         private readonly redisService: RedisService,
//         private readonly dbService: DbService,
//         private readonly pubSubService: PubSubService,
//         private readonly logger: LoggerService
//     ) { }

//     /**
//      * Handles user registration.
//      * - Checks if user exists in Redis or DB.
//      * - Publishes user details for further processing
//      * - add in cache
//      * - Generates a JWT token for the user.
//      * @param email - User's email
//      * @param password - User's password
//      * @returns Registration message with a token.
//      */
//     async register(email: string, password: string) {
//         return this.tracer.startActiveSpan('auth.register', async (span) => {
//             try {
//                 const redis = this.redisService.getClient();
//                 const userCacheKey = `user:${email}`;

//                 if (await redis.exists(userCacheKey)) {
//                     this.logger.warn(`User already exists in Redis: ${email}`);
//                     span.setAttribute('auth.register.status', 'exists-redis');
//                     throw new UnauthorizedException('User already exists');
//                 }

//                 const users = await this.dbService.query('SELECT * FROM users WHERE email = $1', [email]);
//                 if (users.length > 0) {
//                     this.logger.warn(`User already exists in DB: ${email}`);
//                     span.setAttribute('auth.register.status', 'exists-db');
//                     throw new UnauthorizedException('User already exists');
//                 }

//                 const userId = generateUniqueId();

//                 await this.pubSubService.publishUser(email, password, userId);

//                 await redis.set(userCacheKey, JSON.stringify({ email, password, userId }), 'EX', 86400);

//                 const token = this.jwtService.sign({ email, userId });
//                 this.logger.info(`User registered successfully: ${email}`);
//                 span.setAttribute('auth.register.status', 'success');

//                 return { token };
//             } catch (error) {
//                 span.setAttribute('auth.register.error', error.message);
//                 this.logger.error(`Registration failed: ${email}`, { error: error.message });
//                 throw error;
//             } finally {
//                 span.end();
//             }
//         });
//     }

//     /**
//      * Handles user login.
//      * - Checks Redis for cached user data.
//      * - If not found, queries the database. and store in cache
//      * - Validates the password and issues a JWT
//      * @param email - User's email
//      * @param password - User's password
//      * @returns JWT token if successful.
//      */
//     async login(email: string, password: string) {
//         return this.tracer.startActiveSpan('auth.login', async (span) => {
//             try {
//                 const redis = this.redisService.getClient();
//                 const userCacheKey = `user:${email}`;
//                 let user = await redis.get(userCacheKey);

//                 if (!user) {
//                     const users = await this.dbService.query('SELECT * FROM users WHERE email = $1', [email]);
//                     if (users.length === 0) {
//                         this.logger.warn(`Login failed: User not found - ${email}`);
//                         span.setAttribute('auth.login.status', 'failed-not-found');
//                         throw new UnauthorizedException('Invalid credentials');
//                     }
//                     user = users[0];
//                     user = JSON.parse(user);
//                     await redis.set(userCacheKey, JSON.stringify({ user.email, user.password, user.userId }), 'EX', 86400);
//                 }else{
//                     user = JSON.parse(user);
//                 }

//                 const isPasswordValid = await bcrypt.compare(password, user.password);
//                 if (!isPasswordValid) {
//                     this.logger.warn(`Login failed: Incorrect password - ${email}`);
//                     span.setAttribute('auth.login.status', 'failed-password');
//                     throw new UnauthorizedException('Invalid credentials');
//                 }

//                 const token = this.jwtService.sign({ email, userId: user.userId });
//                 this.logger.info(`User logged in successfully: ${email}`);
//                 span.setAttribute('auth.login.status', 'success');
//                 return { token };
//             } catch (error) {
//                 span.setAttribute('auth.login.error', error.message);
//                 this.logger.error(`Login failed: ${email}`, { error: error.message });
//                 throw error;
//             } finally {
//                 span.end();
//             }
//         });
//     }
// }

// import { Injectable, UnauthorizedException } from '@nestjs/common';
// import * as bcrypt from 'bcryptjs';
// import { JwtService } from '@nestjs/jwt';
// import { RedisService } from '../redis/redis.service';
// import { DbService } from '../db/db.service';
// import { PubSubService } from '../pubsub/pubsub.service';
// import { LoggerService } from '../logger/logger.service';
// import { trace } from '@opentelemetry/api';
// import { generateUniqueId } from '../ID';

// @Injectable()
// export class AuthService {
//     private readonly tracer = trace.getTracer('nestjs-auth');

//     constructor(
//         private readonly jwtService: JwtService,
//         private readonly redisService: RedisService,
//         private readonly dbService: DbService,
//         private readonly pubSubService: PubSubService,
//         private readonly logger: LoggerService
//     ) { }



//     /**
//      * Handles user registration.
//      * - Checks if user exists in Redis or DB.
//      * - hash password
//      * - Publishes user details for further processing
//      * - add in cache
//      * - Generates a JWT token for the user.
//      * @param email - User's email
//      * @param password - User's password
//      * @returns Registration message with a token.
//      */
//     async register(email: string, password: string) {
//         return this.tracer.startActiveSpan('auth.register', async (span) => {
//             try {
//                 const redis = this.redisService.getClient();
//                 const userCacheKey = `user:${email}`;

//                 if (await redis.exists(userCacheKey)) {
//                     this.logger.warn(`User already exists in Redis: ${email}`);
//                     throw new UnauthorizedException('User already exists');
//                 }

//                 const users = await this.dbService.query('SELECT * FROM users WHERE email = $1', [email]);
//                 if (users.length > 0) {
//                     this.logger.warn(`User already exists in DB: ${email}`);
//                     throw new UnauthorizedException('User already exists');
//                 }

//                 const hashedPassword = await bcrypt.hash(password, 10);
//                 const userId = generateUniqueId();

//                 await this.pubSubService.publishUser(email, hashedPassword, userId);
//                 await redis.set(userCacheKey, JSON.stringify({ email, hashedPassword, userId }), {
//                     EX: 86400,
//                 });

//                 const token = this.jwtService.sign({ email, userId }, { expiresIn: '7d' });
//                 this.logger.info(`User registered successfully: ${email}`);

//                 return { token };
//             } catch (error) {
//                 this.logger.error(`Registration failed: ${email}`, { error: error.message });
//                 throw error;
//             } finally {
//                 span.end();
//             }
//         });
//     }



//     /**
//      * Handles user login.
//      * - Checks Redis for cached user data.
//      * - If not found, queries the database. and store in cache
//      * - Validates the password and issues a JWT
//      * @param email - User's email
//      * @param password - User's password
//      * @returns JWT token if successful.
//      */
//     async login(email: string, password: string) {
//         return this.tracer.startActiveSpan('auth.login', async (span) => {
//             try {
//                 const redis = this.redisService.getClient();
//                 const userCacheKey = `user:${email}`;
//                 let user: { email: string; password: string; userId: string } | null;

//                 const cachedUser = await redis.get(userCacheKey);

//                 if (cachedUser) {
//                     user = JSON.parse(cachedUser);
//                 } else {
//                     const users = await this.dbService.query('SELECT * FROM users WHERE email = $1', [email]);

//                     if (users.length === 0) {
//                         this.logger.warn(`Login failed: User not found - ${email}`);
//                         throw new UnauthorizedException('Invalid credentials');
//                     }

//                     user = users[0];

//                     if (!user || !user.email || !user.password || !user.userId) {
//                         throw new Error("User data is invalid or incomplete");
//                     }

//                     await redis.set(userCacheKey, JSON.stringify(user), { EX: 86400 });
//                 }

//                 const isPasswordValid = await bcrypt.compare(password, user.password);
//                 if (!isPasswordValid) {
//                     this.logger.warn(`Login failed: Incorrect password - ${email}`);
//                     throw new UnauthorizedException('Invalid credentials');
//                 }

//                 const token = this.jwtService.sign({ email, userId: user.userId }, { expiresIn: '7d' });
//                 this.logger.info(`User logged in successfully: ${email}`);
//                 return { token };
//             } catch (error) {
//                 this.logger.error(`Login failed: ${email}`, { error: error.message });
//                 throw error;
//             } finally {
//                 span.end();
//             }
//         });
//     }

// }

// import { Injectable, UnauthorizedException } from '@nestjs/common';
// import * as bcrypt from 'bcryptjs';
// import { JwtService } from '@nestjs/jwt';
// import { RedisService } from '../redis/redis.service';
// import { DbService } from '../db/db.service';
// import { PubSubService } from '../pubsub/pubsub.service';
// import { LoggerService } from '../logger/logger.service';
// import { trace } from '@opentelemetry/api';
// import { generateUniqueId } from '../ID';

// @Injectable()
// export class AuthService {
//   private readonly tracer = trace.getTracer('nestjs-auth');

//   constructor(
//     private readonly jwtService: JwtService,
//     private readonly redisService: RedisService,
//     private readonly dbService: DbService,
//     private readonly pubSubService: PubSubService,
//     private readonly logger: LoggerService
//   ) { }

//   async register(email: string, password: string) {
//     return this.tracer.startActiveSpan('auth.register', async (span) => {
//       try {
//         const redis = this.redisService.getClient();
//         const userCacheKey = `user:${email}`;

//         if (await redis.exists(userCacheKey)) {
//           this.logger.warn(`User already exists in Redis: ${email}`);
//           throw new UnauthorizedException('User already exists');
//         }

//         const users = await this.dbService.query('SELECT * FROM users WHERE email = $1', [email]);
//         if (users.length > 0) {
//           this.logger.warn(`User already exists in DB: ${email}`);
//           throw new UnauthorizedException('User already exists');
//         }

//         const hashedPassword = await bcrypt.hash(password, 10);
//         const userId = generateUniqueId();

//         await this.pubSubService.publishUser(email, hashedPassword, userId);
//         await redis.set(userCacheKey, JSON.stringify({
//           email, hashedPassword, userId: BigInt(userId).toString(),
//         }), {
//           EX: 86400,
//         });

//         const token = this.jwtService.sign({ email, userId }, { expiresIn: '7d' });
//         this.logger.info(`User registered successfully: ${email}`);

//         return { token };
//       } catch (error) {
//         this.logger.error(`Registration failed: ${email}`, { error: error.message });
//         throw error;
//       } finally {
//         span.end();
//       }
//     });
//   }

//   async login(email: string, password: string) {
//     return this.tracer.startActiveSpan('auth.login', async (span) => {
//       try {
//         const redis = this.redisService.getClient();
//         const userCacheKey = `user:${email}`;
//         let user: { email: string; password: string; userId: string } | null;

//         const cachedUser = await redis.get(userCacheKey);

//         if (cachedUser) {
//           user = JSON.parse(cachedUser);
//         } else {
//           const users = await this.dbService.query('SELECT * FROM users WHERE email = $1', [email]);

//           if (users.length === 0) {
//             this.logger.warn(`Login failed: User not found - ${email}`);
//             throw new UnauthorizedException('Invalid credentials');
//           }

//           user = users[0];

//           if (!user || !user.email || !user.password || !user.userId) {
//             throw new Error("User data is invalid or incomplete");
//           }

//           await redis.set(userCacheKey, JSON.stringify({user.email, user.password, userId: BigInt(userId).toString()}), { EX: 86400 });
//         }

//         const isPasswordValid = await bcrypt.compare(password, user.password);
//         if (!isPasswordValid) {
//           this.logger.warn(`Login failed: Incorrect password - ${email}`);
//           throw new UnauthorizedException('Invalid credentials');
//         }

//         const token = this.jwtService.sign({ email, userId: Number(user.userId) }, { expiresIn: '7d' });
//         this.logger.info(`User logged in successfully: ${email}`);
//         return { token };
//       } catch (error) {
//         this.logger.error(`Login failed: ${email}`, { error: error.message });
//         throw error;
//       } finally {
//         span.end();
//       }
//     });
//   }
// }


import { Injectable, UnauthorizedException } from '@nestjs/common';
import * as bcrypt from 'bcryptjs';
import { JwtService } from '@nestjs/jwt';
import { RedisService } from '../redis/redis.service';
import { DbService } from '../db/db.service';
import { PubSubService } from '../pubsub/pubsub.service';
import { LoggerService } from '../logger/logger.service';
import { trace } from '@opentelemetry/api';
import { generateUniqueId } from '../ID';

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

  async register(email: string, password: string) {
    return this.tracer.startActiveSpan('auth.register', async (span) => {
      try {
        const redis = this.redisService.getClient();
        const userCacheKey = `user:${email}`;

        if (await redis.exists(userCacheKey)) {
          this.logger.warn(`User already exists in Redis: ${email}`);
          throw new UnauthorizedException('User already exists');
        }

        const users = await this.dbService.query('SELECT * FROM users WHERE email = $1', [email]);
        if (users.length > 0) {
          this.logger.warn(`User already exists in DB: ${email}`);
          throw new UnauthorizedException('User already exists');
        }

        const hashedPassword = await bcrypt.hash(password, 10);
        const userId = generateUniqueId(); // Ensure this returns a BIGINT number

        // Store in PubSub (Assuming it handles BIGINT properly)
        await this.pubSubService.publishUser(email, hashedPassword, userId);

        // Store in Redis (Convert userId to string for compatibility)
        await redis.set(userCacheKey, JSON.stringify({
          email,
          password: hashedPassword,
          userId: userId.toString(), // Convert to string
        }), { EX: 86400 });

        // Store in DB
        await this.dbService.query(
          'INSERT INTO users (id, email, password) VALUES ($1, $2, $3)',
          [userId, email, hashedPassword]
        );

        // JWT does not support BigInt, so convert it to Number
        const token = this.jwtService.sign({ email, userId: Number(userId) }, { expiresIn: '7d' });
        this.logger.info(`User registered successfully: ${email}`);

        return { token };
      } catch (error) {
        this.logger.error(`Registration failed: ${email}`, { error: error.message });
        throw error;
      } finally {
        span.end();
      }
    });
  }

  // async login(email: string, password: string) {
  //   return this.tracer.startActiveSpan('auth.login', async (span) => {
  //     try {
  //       const redis = this.redisService.getClient();
  //       const userCacheKey = `user:${email}`;
  //       let user: { email: string; password: string; userId: string } | null;

  //       const cachedUser = await redis.get(userCacheKey);

  //       if (cachedUser) {
  //         user = JSON.parse(cachedUser);
  //         // Convert userId back to BigInt
  //         user.userId = BigInt(user.userId).toString();
  //       } else {
  //         const users = await this.dbService.query('SELECT * FROM users WHERE email = $1', [email]);

  //         if (users.length === 0) {
  //           this.logger.warn(`Login failed: User not found - ${email}`);
  //           throw new UnauthorizedException('Invalid credentials');
  //         }

  //         user = users[0];

  //         if (!user || !user.email || !user.password || !user.userId) {
  //           throw new Error("User data is invalid or incomplete");
  //         }

  //         // Store in Redis for faster future lookups
  //         await redis.set(userCacheKey, JSON.stringify({
  //           email: user.email,
  //           password: user.password,
  //           userId: BigInt(user.userId).toString(),
  //         }), { EX: 86400 });
  //       }

  //       const isPasswordValid = await bcrypt.compare(password, user.password);
  //       if (!isPasswordValid) {
  //         this.logger.warn(`Login failed: Incorrect password - ${email}`);
  //         throw new UnauthorizedException('Invalid credentials');
  //       }

  //       // Convert userId to Number for JWT
  //       const token = this.jwtService.sign({ email, userId: Number(user.userId) }, { expiresIn: '7d' });
  //       this.logger.info(`User logged in successfully: ${email}`);
  //       return { token };
  //     } catch (error) {
  //       this.logger.error(`Login failed: ${email}`, { error: error.message });
  //       throw error;
  //     } finally {
  //       span.end();
  //     }
  //   });
  // }
  async login(email: string, password: string) {
    return this.tracer.startActiveSpan('auth.login', async (span) => {
      try {
        const redis = this.redisService.getClient();
        const userCacheKey = `user:${email}`;
        let user: { email: string; password: string; userId: string } | null;
  
        const cachedUser = await redis.get(userCacheKey);
  
        console.log(cachedUser);
        if (cachedUser) {
          user = JSON.parse(cachedUser);
  
          // Ensure userId is stored as a valid BigInt string
          if (!user.userId) {
            this.logger.error(`Login failed: Missing userId in Redis for ${email}`);
            throw new UnauthorizedException('Invalid credentials');
          }
  
          user.userId = BigInt(user.userId).toString();
          console.log(user.userId);
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
  
        // Validate password
        console.log(password, user.password)
        const isPasswordValid = await bcrypt.compare(password, user.password);
        if (!isPasswordValid) {
          this.logger.warn(`Login failed: Incorrect password for ${email}`);
          throw new UnauthorizedException('Invalid credentials');
        }
        console.log(isPasswordValid);
  
        // Convert userId to Number for JWT
        const token = this.jwtService.sign(
          { email, userId: Number(user.userId) },
          { expiresIn: '7d' }
        );
  
        console.log(token)
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
