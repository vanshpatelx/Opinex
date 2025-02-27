// src/config/DB/query.ts
/**
    PostgreSQL Queries (User Database)

    Contains SQL query definitions for interacting with the `Auth` table.

    Queries:
    - `addUser` → Inserts a new user record into the `Auth` table.

    Dependencies:
    - PostgreSQL client for executing queries

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: February 2, 2025  
 */

const queryList = {
    addUser: `INSERT INTO "Auth" (id, email, password, type, created_at, updated_at, inserted_at) VALUES ($1, $2, $3, $4, $5, $6, $7);`,
};

export default queryList;
