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

package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var (
	dbPool *sql.DB
	ctx    = context.Background()
)

// Init initializes the database connection pool
func Init(databaseURL string) error {
	var err error
	dbPool, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	dbPool.SetMaxOpenConns(20)
	dbPool.SetMaxIdleConns(10)
	dbPool.SetConnMaxLifetime(5 * time.Minute)

	for attempt := 1; attempt <= 5; attempt++ {
		if err = dbPool.Ping(); err == nil {
			return nil
		}
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("could not connect to PostgreSQL after multiple attempts: %w", err)
}

// FetchOne executes a query and returns a single row result
func FetchOne(query string, args ...interface{}) (map[string]interface{}, error) {
	rows, err := dbPool.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names
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

	// Scan row values
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
func FetchAll(query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := dbPool.QueryContext(ctx, query, args...)
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

// Close database connection pool
func Close() {
	if dbPool != nil {
		_ = dbPool.Close()
	}
}
