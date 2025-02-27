// src/routes/auth.ts
/**
    Authentication Routes

    Endpoints:
    - POST /auth/login → User login, returns JWT.  
    - POST /auth/register → User registration.

    Dependencies:
    - Express Router  
    - Controllers: loginUser, registerUser  

    Author: Vansh Patel (remotevansh@gmail.com)  
    Date: February 2, 2025  
 */


import { Router } from "express";
import { loginUser, registerUser } from "../controllers/auth";

const router = Router();

router.post("/login", loginUser);
router.post("/register", registerUser);

export default router;
