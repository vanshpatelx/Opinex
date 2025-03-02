import { describe, it, expect, beforeAll, test } from 'vitest';
import axios from 'axios';

const AUTH_URL = "http://localhost:5001/auth";
const HOLDING_URL = "http://localhost:8081/holdings";

const tempUser = {
    email: "testHolding1@example.com",
    password: "Test@1234"
};

const normalUser = {
    email: "testuser@gmail.com",
    password: "12345678"
}
const adminUser = {
    email: "test@gmail.com",
    password: "12345678"
};

let adminToken = "";
let tempToken = "";
let normalToken = "";

beforeAll(async () => {
    const tempResponse = await axios.post(`${AUTH_URL}/register`, tempUser);
    expect(tempResponse.status).toBe(201);
    tempToken = tempResponse.data.token;

    const adminResponse = await axios.post(`${AUTH_URL}/login`, adminUser);
    expect(adminResponse.status).toBe(200);
    adminToken = adminResponse.data.token;

    const normalResponse = await axios.post(`${AUTH_URL}/login`, normalUser);
    expect(normalResponse.status).toBe(200);
    normalToken = normalResponse.data.token;
});

describe("Get Balance", () => {
    it("Get Balance normal User", async () => {
        const response = await axios.get(`${HOLDING_URL}/balance`, {
            params: { userID: 2 },
            headers: { Authorization: `Bearer ${normalToken}` }
        });
        expect(response.status).toBe(200);
        expect(response.data).toHaveProperty("balance");
        expect(response.data).toHaveProperty("locked_balance");
    });

    it("Check user only access their own balance means try by another user", async () => {
        try {
            await axios.get(`${HOLDING_URL}/balance`, {
                params: { userID: 2 },
                headers: { Authorization: `Bearer ${tempToken}` }
            });
        } catch (error: any) {
            expect(error.response?.status).toBe(401);
        }
    });

    it("Admin can access anyone balance", async () => {
        const response = await axios.get(`${HOLDING_URL}/balance`, {
            params: { userID: 2 },
            headers: { Authorization: `Bearer ${adminToken}` }
        });
        expect(response.status).toBe(200);
        expect(response.data).toHaveProperty("balance");
        expect(response.data).toHaveProperty("locked_balance");
    });

    it("Without token", async () => {
        try {
            await axios.get(`${HOLDING_URL}/balance`, {
                params: { userID: 2 },
            });
        } catch (error: any) {
            expect(error.response?.status).toBe(401);
        }
    });
});
