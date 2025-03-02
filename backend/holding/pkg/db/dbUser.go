/***
Package db - PostgreSQL Database Connection for User Service

Handles the database connection pooling, query execution, and graceful shutdown 
for the User Service.

Features:
1. **Singleton Connection Pool** - Ensures a single database connection instance.
2. **Connection Management** - Configures max open/idle connections and retries on failure.
3. **Query Execution Helpers** - Fetches single or multiple records safely.
4. **Graceful Shutdown Support** - Closes the connection when shutting down.

Methods:
- **InitDBUser**: Initializes the PostgreSQL connection with retry logic.
- **GetDBUser**: Returns the singleton database connection.
- **FetchOneUser**: Retrieves a single row from the database.
- **FetchAllUser**: Retrieves multiple rows from the database.
- **CloseDBUser**: Closes the database connection.

Last Updated: March 2, 2025
**/


package db

import (
	"context"
	"database/sql"
	"fmt"
	"holding/pkg/config"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	onceDBUser  sync.Once
	dbPoolUser  *sql.DB
	ctxDBUser   = context.Background()
	initErrUser error
)

// InitDB initializes the database connection as a singleton
func InitDBUser(databaseURL string) error {
	onceDBUser.Do(func() {
		dbPoolUser, initErrUser = sql.Open("postgres", databaseURL)
		if initErrUser != nil {
			initErrUser = fmt.Errorf("failed to open database connection: %w", initErrUser)
			return
		}

		dbPoolUser.SetMaxOpenConns(20)
		dbPoolUser.SetMaxIdleConns(10)
		dbPoolUser.SetConnMaxLifetime(5 * time.Minute)

		// Retry connection with exponential backoff
		for attempt := 1; attempt <= 5; attempt++ {
			if err := dbPoolUser.Ping(); err == nil {
				initErrUser = nil
				return
			}
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
		}

		initErrUser = fmt.Errorf("could not connect to PostgreSQL after multiple attempts")
	})

	return initErrUser
}

// GetDB returns the singleton database connection
func GetDBUser() (*sql.DB, error) {
	if dbPoolUser == nil {
		if err := InitDBUser(config.AppConfig.USER_DBURL); err != nil {
			return nil, err
		}
	}
	return dbPoolUser, nil
}

// FetchOne executes a query and returns a single row result
func FetchOneUser(query string, args ...interface{}) (map[string]interface{}, error) {
	dbUser, err := GetDBUser()
	if err != nil {
		return nil, err
	}

	rows, err := dbUser.QueryContext(ctxDBUser, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("no rows found")
	}

	values := make([]interface{}, len(colTypes))
	valuePtrs := make([]interface{}, len(colTypes))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for i, col := range colTypes {
		result[col.Name()] = values[i]
	}

	return result, nil
}

// FetchAll executes a query and returns multiple rows
func FetchAllUser(query string, args ...interface{}) ([]map[string]interface{}, error) {
	dbUser, err := GetDBUser()
	if err != nil {
		return nil, err
	}

	rows, err := dbUser.QueryContext(ctxDBUser, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(colTypes))
		valuePtrs := make([]interface{}, len(colTypes))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, col := range colTypes {
			rowMap[col.Name()] = values[i]
		}

		results = append(results, rowMap)
	}

	return results, nil
}

// CloseDB closes the database connection
func CloseDBUser() {
	if dbPoolUser != nil {
		_ = dbPoolUser.Close()
	}
}
