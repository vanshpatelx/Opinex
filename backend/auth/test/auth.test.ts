import request from "supertest";
import { describe, it, expect, beforeAll, afterAll } from "vitest";
import app from "../src/index";
import { redisClient } from "../src/config/Cache/RedisClient";
// import { postgresClient } from "../src/config/DB/db";

const AUTH_URL = "/auth";

describe("Auth Service", () => {
    const testUser = {
        email: `testuser@example.com`,
        password: "Test@1234"
    };

    beforeAll(async () => {
        const cachedUser = await redisClient.del(`user:${testUser.email}`);
        expect(cachedUser).toBe(1);

    })

    let token: string;

    it("should register a new user", async () => {
        const response = await request(app)
            .post(`${AUTH_URL}/register`)
            .send(testUser);

        expect(response.status).toBe(201);
        expect(response.body).toHaveProperty("message", "User registered successfully");
        expect(response.body).toHaveProperty("token");

        token = response.body.token;
    });

    // it("should remove user data from Redis cache after registration", async () => {
    //     const cachedUser = await redisClient.del(`user:${testUser.email}`);
    //     expect(cachedUser).toBe(1);
    // });

    // it("should not register the same user again", async () => {
    //     const response = await request(app)
    //         .post(`${AUTH_URL}/register`)
    //         .send(testUser);

    //     expect(response.status).toBe(400);
    //     expect(response.body).toHaveProperty("message", "User already exists");
    // });

    it("should log in an existing user", async () => {
        const response = await request(app)
            .post(`${AUTH_URL}/login`)
            .send(testUser);

        expect(response.status).toBe(200);
        expect(response.body).toHaveProperty("message", "User login successful");
        expect(response.body).toHaveProperty("token");
    });

    it("should not log in with incorrect password", async () => {
        const response = await request(app)
            .post(`${AUTH_URL}/login`)
            .send({ ...testUser, password: "WrongPassword" });

        expect(response.status).toBe(401);
        expect(response.body).toHaveProperty("message", "Invalid credentials");
    });

    it("should return error when logging in with missing fields", async () => {
        const response = await request(app)
            .post(`${AUTH_URL}/login`)
            .send({ email: testUser.email });

        expect(response.status).toBe(500);
    });

    // it("remove user session from Redis cache", async () => {
    //     const cachedUser = await redisClient.del(`user:${testUser.email}`);
    //     expect(cachedUser).toBe(1);
    // });

    // it("should still log in when Redis is empty", async () => {
    //     const response = await request(app)
    //         .post(`${AUTH_URL}/login`)
    //         .send(testUser);

    //     expect(response.status).toBe(200);
    //     expect(response.body).toHaveProperty("message", "User login successful");
    //     expect(response.body).toHaveProperty("token");
    // });
});
