package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"context"
)

var db *sql.DB
var rdb *redis.Client
var ctx = context.Background()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Database and Redis
	initializeDatabase()

	// Set up routes
	r := mux.NewRouter()
	r.HandleFunc("/users", addUser).Methods("POST")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// Start the server
	http.Handle("/", r)
	log.Printf("Server running on port %s", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func initializeDatabase() {
	var err error

	// Initialize PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	// Initialize Redis
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})

	// Create users table if it does not exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL
		);
	`)
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}

	log.Println("Database and Redis initialized successfully.")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "message": "Server is running."}`))
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var user struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Parse JSON body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert user into PostgreSQL
	query := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3) RETURNING id, name, email`
	row := db.QueryRow(query, user.ID, user.Name, user.Email)

	// Check for any errors
	err = row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		http.Error(w, "Error inserting user into DB", http.StatusInternalServerError)
		return
	}

	// Store user in Redis
	userData, _ := json.Marshal(user)
	err = rdb.Set(ctx, fmt.Sprintf("user:%d", user.ID), userData, 0).Err()
	if err != nil {
		http.Error(w, "Error saving user to Redis", http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "message": "User added successfully."}`))
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Try to get user from Redis
	redisData, err := rdb.Get(ctx, "user:"+id).Result()
	if err == redis.Nil {
		// If not found in Redis, get from DB
		var user struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		query := `SELECT id, name, email FROM users WHERE id = $1`
		row := db.QueryRow(query, id)
		err := row.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Store user in Redis
		userData, _ := json.Marshal(user)
		err = rdb.Set(ctx, fmt.Sprintf("user:%d", user.ID), userData, 0).Err()
		if err != nil {
			http.Error(w, "Error saving user to Redis", http.StatusInternalServerError)
			return
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"success": true, "user": %s}`, redisData)))
	} else if err != nil {
		http.Error(w, "Error retrieving user from Redis", http.StatusInternalServerError)
		return
	} else {
		// If found in Redis, return the cached data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"success": true, "user": %s}`, redisData)))
	}
}