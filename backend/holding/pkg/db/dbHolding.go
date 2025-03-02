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
	"holding/pkg/config"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	onceDBHolding  sync.Once
	dbPoolHolding  *sql.DB
	ctxDBHolding   = context.Background()
	initErrHolding error
)

// InitDB initializes the database connection as a singleton
func InitDBHolding(databaseURL string) error {
	onceDBHolding.Do(func() {
		dbPoolHolding, initErrHolding = sql.Open("postgres", databaseURL)
		if initErrHolding != nil {
			initErrHolding = fmt.Errorf("failed to open database connection: %w", initErrHolding)
			return
		}

		dbPoolHolding.SetMaxOpenConns(20)
		dbPoolHolding.SetMaxIdleConns(10)
		dbPoolHolding.SetConnMaxLifetime(5 * time.Minute)

		// Retry connection with exponential backoff
		for attempt := 1; attempt <= 5; attempt++ {
			if err := dbPoolHolding.Ping(); err == nil {
				initErrHolding = nil
				return
			}
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
		}

		initErrHolding = fmt.Errorf("could not connect to PostgreSQL after multiple attempts")
	})

	return initErrHolding
}

// GetDB returns the singleton database connection
func GetDBHolding() (*sql.DB, error) {
	if dbPoolHolding == nil {
		if err := InitDBHolding(config.AppConfig.HOLDING_DBURL); err != nil {
			return nil, err
		}
	}
	return dbPoolHolding, nil
}

// FetchOne executes a query and returns a single row result
func FetchOneHolding(query string, args ...interface{}) (map[string]interface{}, error) {
	dbHolding, err := GetDBHolding()
	if err != nil {
		return nil, err
	}

	rows, err := dbHolding.QueryContext(ctxDBHolding, query, args...)
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
func FetchAllHolding(query string, args ...interface{}) ([]map[string]interface{}, error) {
	dbHolding, err := GetDBHolding()
	if err != nil {
		return nil, err
	}

	rows, err := dbHolding.QueryContext(ctxDBHolding, query, args...)
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
func CloseDBHolding() {
	if dbPoolHolding != nil {
		_ = dbPoolHolding.Close()
	}
}
