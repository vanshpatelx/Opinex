// src/config/DB/query.ts
/**
    DB Queries

    Stores SQL queries related to all.

    Queries:
    - getUser → Retrieves user details (id, email, password, type) by email.

    Dependencies:
    - PostgreSQL

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: February 2, 2025  
 */

const queryList = {
    getUser: `SELECT id, email, password, type FROM "Auth" WHERE email = $1`
};

export default queryList;
