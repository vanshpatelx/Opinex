from config.config import config

def load_env():
    print("==================================")
    print("🚀 Application Environment Variables")
    print("==================================\n")

    print("🔹 REDIS CONFIGURATION:")
    print("   ➤ REDIS_HOST      :", config.REDIS_HOST)
    print("   ➤ REDIS_PORT      :", config.REDIS_PORT)
    print("   ➤ REDIS_PASSWORD  :", config.REDIS_PASSWORD)
    print("\n")

    print("🔹 JWT CONFIGURATION:")
    print("   ➤ JWT_SECRET      :", config.JWT_SECRET)
    print("\n")

    print("🔹 DATABASE CONFIGURATION:")
    print("   ➤ DB_USER         :", config.DB_USER)
    print("   ➤ DB_HOST         :", config.DB_HOST)
    print("   ➤ DB_NAME         :", config.DB_NAME)
    print("   ➤ DB_PASSWORD     :", config.DB_PASSWORD)
    print("   ➤ DB_PORT         :", config.DB_PORT)
    print("\n")

    print("🔹 SERVER PORT:")
    print("   ➤ PORT            :", config.PORT)
    print("\n")

    print("🔹 RABBITMQ CONFIGURATION:")
    print("   ➤ RABBITMQ_HOST   :", config.RABBITMQ_HOST)
    print("   ➤ RABBITMQ_PORT   :", config.RABBITMQ_PORT)
    print("   ➤ RABBITMQ_USER   :", config.RABBITMQ_USER)
    print("   ➤ RABBITMQ_PASSWORD:", config.RABBITMQ_PASSWORD)
    print("   ➤ RABBITMQ_EXCHANGES:", config.RABBITMQ_EXCHANGES)
    print("   ➤ RABBITMQ_URL    :", config.RABBITMQ_URL)
    print("\n")

    print("==================================")
    print("✅ Environment variables loaded!")
    print("==================================")

# Run the loader
if __name__ == "__main__":
    load_env()
