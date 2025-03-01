import { describe, it, expect } from 'vitest';
import axios from 'axios';

const AUTH_URL = "http://localhost:5001/auth";

describe("Auth Service", () => {

    const testUser = {
        email: `test974@example.com`,
        password: "Test@1234"
    };

    it("Register a user", async () => {
        const response = await axios.post(`${AUTH_URL}/register`, testUser);
        
        expect(response.status).toBe(201);
        expect(response.data).toHaveProperty("message", "User registered successfully");
        expect(response.data).toHaveProperty("token");        
    });

    it("Try to register the same user", async () => {
        try {
            await axios.post(`${AUTH_URL}/register`, testUser);
        } catch (error) {
            expect(error.response.status).toBe(400);
            expect(error.response.data).toHaveProperty("message", "User already exists");
        }
    });

    it("Try to register a user with invalid email", async () => {
        const testUser = {
            email: `invalid@.com`,
            password: "Test@1234"
        };

        try {
            await axios.post(`${AUTH_URL}/register`, testUser);
        } catch (error) {
            expect(error.response.status).toBe(500);
            expect(error.response.data).toHaveProperty("message", "Internal server error");
        }
    });

    it("Try to register a user with invalid password", async () => {
        const testUser = {
            email: `valid@gmail.com`,
            password: '',
        };

        try {
            await axios.post(`${AUTH_URL}/register`, testUser);
        } catch (error) {
            expect(error.response.status).toBe(500); 
            expect(error.response.data).toHaveProperty("message", "Internal server error");      
        }
    });

    it("Try to register a user with invalid Params or Body", async () => {
        const testUser = {
            email: `valid@gmail.com`,
        };

        try {
            await axios.post(`${AUTH_URL}/register`, testUser);
        } catch (error) {
            expect(error.response.status).toBe(500);
            expect(error.response.data).toHaveProperty("message", "Internal server error");
        }
    });


    it("Login an existing user", async () => {
        const response = await axios.post(`${AUTH_URL}/login`, testUser);
        
        expect(response.status).toBe(200);
        expect(response.data).toHaveProperty("message", "User login successfully");
        expect(response.data).toHaveProperty("token");
    });


    it("Try to login a user with invalid email", async () => {
        const testUser = {
            email: `invalid@.com`,
            password: "Test@1234"
        };

        try {
            await axios.post(`${AUTH_URL}/login`, testUser);
        } catch (error) {
            expect(error.response.status).toBe(500);
            expect(error.response.data).toHaveProperty("message", "Internal server error");
        }
    });


    it("Try to login a user with not registered email", async () => {
        const testUser = {
            email: `testInvalid@gmail.com`,
            password: "Test@1234"
        };
    
        const response = await axios.post(`${AUTH_URL}/login`, testUser, {
            validateStatus: () => true // Prevents Axios from throwing errors
        });
    
        expect(response.status).toBe(401);
        expect(response.data).toHaveProperty("message", "Invalid credentials");
    });
    

    it("Try to login a user with invalid password", async () => {
        try {
            await axios.post(`${AUTH_URL}/login`, { ...testUser, password: "WrongPassword" });
        } catch (error) {
            expect(error.response.status).toBe(401);
            expect(error.response.data).toHaveProperty("message", "Invalid credentials");
        }
    });

    it("Try to login a user with invalid Params or Body", async () => {
        const testUser = {
            email: `testInvalid@gmail.com`,
        };
        try {
            await axios.post(`${AUTH_URL}/login`, testUser);
        } catch (error) {
            expect(error.response.status).toBe(500);
            expect(error.response.data).toHaveProperty("message", "Internal server error");
        }
    });
});
