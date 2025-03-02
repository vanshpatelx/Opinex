// holding/pkg/db/db.go
/***
PostgreSQL Database Connection

Manages connection pooling and queries for PostgreSQL.

Features:
1. Initializes a connection pool to PostgreSQL.
2. Provides methods for executing queries safely.
3. Supports graceful shutdown.

Author: Vansh Patel (remotevansh@gmail.com)
Last Updated: February 26, 2025
***/

// package db

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"time"

// 	_ "github.com/lib/pq"
// )

// var (
// 	dbPool *sql.DB
// 	ctx    = context.Background()
// )

// // Init initializes the database connection pool
// func Init(databaseURL string) error {
// 	var err error
// 	dbPool, err = sql.Open("postgres", databaseURL)
// 	if err != nil {
// 		return fmt.Errorf("failed to open database connection: %w", err)
// 	}

// 	dbPool.SetMaxOpenConns(20)
// 	dbPool.SetMaxIdleConns(10)
// 	dbPool.SetConnMaxLifetime(5 * time.Minute)

// 	for attempt := 1; attempt <= 5; attempt++ {
// 		if err = dbPool.Ping(); err == nil {
// 			return nil
// 		}
// 		time.Sleep(2 * time.Second)
// 	}

// 	return fmt.Errorf("could not connect to PostgreSQL after multiple attempts: %w", err)
// }

// // FetchOne executes a query and returns a single row result
// func FetchOne(query string, args ...interface{}) (map[string]interface{}, error) {
// 	rows, err := dbPool.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	// Get column names
// 	colTypes, err := rows.ColumnTypes()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if !rows.Next() {
// 		return nil, fmt.Errorf("no rows found")
// 	}

// 	values := make([]interface{}, len(colTypes))
// 	valuePtrs := make([]interface{}, len(colTypes))
// 	for i := range values {
// 		valuePtrs[i] = &values[i]
// 	}

// 	// Scan row values
// 	if err := rows.Scan(valuePtrs...); err != nil {
// 		return nil, err
// 	}

// 	result := make(map[string]interface{})
// 	for i, col := range colTypes {
// 		result[col.Name()] = values[i]
// 	}

// 	return result, nil
// }

// // FetchAll executes a query and returns multiple rows
// func FetchAll(query string, args ...interface{}) ([]map[string]interface{}, error) {
// 	rows, err := dbPool.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	colTypes, err := rows.ColumnTypes()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var results []map[string]interface{}

// 	for rows.Next() {
// 		values := make([]interface{}, len(colTypes))
// 		valuePtrs := make([]interface{}, len(colTypes))
// 		for i := range values {
// 			valuePtrs[i] = &values[i]
// 		}

// 		if err := rows.Scan(valuePtrs...); err != nil {
// 			return nil, err
// 		}

// 		rowMap := make(map[string]interface{})
// 		for i, col := range colTypes {
// 			rowMap[col.Name()] = values[i]
// 		}

// 		results = append(results, rowMap)
// 	}

// 	return results, nil
// }

// // Close database connection pool
// func Close() {
// 	if dbPool != nil {
// 		_ = dbPool.Close()
// 	}
// }



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
		if err := InitDBUser(config.AppConfig.HOLDING_DBURL); err != nil {
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
