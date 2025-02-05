import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import './otel.tracing';

async function bootstrap() {
  printEnvVariables();
  const app = await NestFactory.create(AppModule);
  const port = Number(process.env.AUTH_PORT);
  await app.listen(port);
  console.log(`🚀 NestJS: Auth server is running on ${port}`);
}

bootstrap();


function printEnvVariables() {
  const colors = {
    reset: '\x1b[0m',
    cyan: '\x1b[36m',
    yellow: '\x1b[33m',
    green: '\x1b[32m',
    red: '\x1b[31m',
    blue: '\x1b[34m',
    magenta: '\x1b[35m',
    bold: '\x1b[1m',
  };

  console.log(
    `\n${colors.bold}${colors.cyan}=======================================${colors.reset}`
  );
  console.log(
    `${colors.bold}${colors.green}🚀 NestJS: Auth Server - Environment Variables${colors.reset}`
  );
  console.log(
    `${colors.bold}${colors.cyan}=======================================${colors.reset}\n`
  );

  const envVariables = {
    '📦 Database': {
      'AUTH_DB_HOST': process.env.AUTH_DB_HOST,
      'AUTH_DB_PORT': process.env.AUTH_DB_PORT,
      'AUTH_DB_USER': process.env.AUTH_DB_USER,
      'AUTH_DB_PASSWORD': process.env.AUTH_DB_PASSWORD,
      'AUTH_DB_NAME': process.env.AUTH_DB_NAME,
    },
    '🚀 Redis': {
      'AUTH_REDIS_HOST': process.env.AUTH_REDIS_HOST,
      'AUTH_REDIS_PORT': process.env.AUTH_REDIS_PORT,
    },
    '🖥️ Server': {
      'AUTH_PORT': process.env.AUTH_PORT,
    }
  };

  Object.entries(envVariables).forEach(([section, vars]) => {
    console.log(
      `${colors.bold}${colors.blue}${section}${colors.reset}`
    );
    console.log(
      `${colors.magenta}---------------------------------------${colors.reset}`
    );
    Object.entries(vars).forEach(([key, value]) => {
      console.log(
        `${colors.yellow}${key}:${colors.reset} ${value ? colors.green + value + colors.reset : colors.red + '❌ NOT SET' + colors.reset}`
      );
    });
    console.log('');
  });

  console.log(
    `${colors.bold}${colors.cyan}=======================================${colors.reset}\n`
  );
}