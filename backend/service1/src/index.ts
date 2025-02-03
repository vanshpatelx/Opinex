import { Client } from "pg";
import { createClient } from "redis";
import { config } from "dotenv";
import express from "express";

config();

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware to parse JSON body
app.use(express.json());

// Initialize PostgreSQL Connection
const dbClient = new Client({
  host: process.env.DB_HOST,
  port: Number(process.env.DB_PORT),
  user: process.env.DB_USER,
  password: String(process.env.DB_PASSWORD),
  database: process.env.DB_NAME,
});

// Initialize Redis Client
const redisClient = createClient({
  url: `redis://${process.env.REDIS_HOST}:${process.env.REDIS_PORT}`,
});

const initializeDatabase = async () => {
  try {
    await dbClient.connect();
    await redisClient.connect();
    console.log("Database and Redis initialized successfully.");
  } catch (error) {
    console.error("Error initializing database or Redis:", error);
    process.exit(1);
  }
};

// Add User to Both DB and Redis
app.post("/users", async (req, res) => {
  const { id, name, email } = req.body;
  
  try {
    const query = "INSERT INTO users (id, name, email) VALUES ($1, $2, $3) RETURNING *";
    const values = [id, name, email];
    const result = await dbClient.query(query, values);
    const user = result.rows[0];
    
    await redisClient.set(`user:${user.id}`, JSON.stringify(user));
    res.status(200).json({ success: true, message: "User added successfully." });
  } catch (error) {
    res.status(500).json({ success: false, message: "Error adding user.", error });
  }
});

app.get("/health", async (req, res) => {
  res.status(200).json({ success: true, message: "server is running." });
});

// Retrieve User from Both DB and Redis
app.get("/users/:id", async (req, res) => {
  const { id } = req.params;
  
  try {
    const redisData = await redisClient.get(`user:${id}`);
    if (!redisData) {
      const query = "SELECT * FROM users WHERE id = $1";
      const result = await dbClient.query(query, [id]);
      const user = result.rows[0];
      if (!user) {
        res.status(404).json({ success: false, message: "User not found." });
        return
      }
      await redisClient.set(`user:${id}`, JSON.stringify(user));
      res.status(200).json({ success: true, user });
      return
    }
    res.status(200).json({ success: true, user: JSON.parse(redisData)});
  } catch (error) {
    res.status(500).json({ success: false, message: "Error retrieving user.", error });
  }
});

// Initialize DB and start server
if (require.main === module) {
  initializeDatabase().then(() => {
    app.listen(PORT, () => {
      console.log(`Server running on port ${PORT}`);
    });
  });
}

export default app;
