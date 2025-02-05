import { Injectable, OnModuleInit } from '@nestjs/common';
import { Pool } from 'pg';

@Injectable()
export class DbService implements OnModuleInit {
  private pool: Pool;

  constructor() {
    this.pool = new Pool({
      user: process.env.AUTH_DB_USER,
      host: process.env.AUTH_DB_HOST,
      database: process.env.AUTH_DB_NAME,
      password: process.env.AUTH_DB_PASSWORD,
      port: Number(process.env.AUTH_DB_PORT),
    });
  }

  async onModuleInit() {
    await this.checkConnection();
  }

  private async checkConnection() {
    try {
      const client = await this.pool.connect();
      await client.query('SELECT 1');
      client.release();
      console.log('Database connected successfully');
    } catch (error) {
      console.error('Failed to connect to the database', error);
      throw new Error('Database connection failed');
    }
  }

  async query(queryText: string, values: any[]): Promise<any[]> {
    const result = await this.pool.query(queryText, values);
    return result.rows;
  }
}
