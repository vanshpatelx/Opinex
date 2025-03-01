import { describe, it, expect, beforeAll, test } from 'vitest';
import axios from 'axios';

const AUTH_URL = "http://localhost:5001/auth";
const EVENT_URL = "http://localhost:8000/event";



const normalUser = {
    email: "testEvent26@example.com",
    password: "Test@1234"
};

const adminUser = {
    email: "test@gmail.com",
    password: "12345678"
};

let adminToken = "";
let userToken = "";

// Login for both users and store tokens
beforeAll(async () => {
    const userResponse = await axios.post(`${AUTH_URL}/register`, normalUser);
    expect(userResponse.status).toBe(201);
    userToken = userResponse.data.token;

    const adminResponse = await axios.post(`${AUTH_URL}/login`, adminUser);
    expect(adminResponse.status).toBe(200);
    adminToken = adminResponse.data.token;
});

describe("Create an event", () => {
    const testEvent = {
        name: "hello",
        details: "Event 1 details",
        settlement_time: "2025-02-28 11:11:27"
    };
    
    describe("Event Creation Tests", () => {
        test.each([...Array(10).keys()])("Create event #%d", async (i) => {
            const response = await axios.post(EVENT_URL, testEvent, {
                headers: { Authorization: `Bearer ${adminToken}` }
            });
    
            expect(response.status).toBe(200);
            expect(response.data).toHaveProperty("message", "Event created successfully");
            expect(response.data).toHaveProperty("event_id");
        });
    });
    

    it("Try to create an event with invalid Params or Body", async () => {
        const invalidEvent = {
            name: "hello",
            details: "Event 1 details",
        };
        try {
            await axios.post(`${EVENT_URL}`, invalidEvent, {
                headers: { Authorization: `Bearer ${adminToken}` }
            });
        } catch (error: any) {
            expect(error.response?.status).toBe(422);
        }
    });

    it("Unauthorized: Normal user tries to create an event", async () => {
        try {
            await axios.post(`${EVENT_URL}`, testEvent, {
                headers: { Authorization: `Bearer ${userToken}` }
            });
        } catch (error: any) {
            expect(error.response?.status).toBe(401);
        }
    });

    it("Without token: Try to create an event", async () => {
        try {
            await axios.post(`${EVENT_URL}`, testEvent);
        } catch (error: any) {
            expect(error.response?.status).toBe(403);
        }
    });
});

describe("Get an event", () => {
    const testEvent = {
        name: "hello",
        details: "Event 1 details",
        settlement_time: "2025-02-28 11:11:27"
    };

    let eventID: bigint;

    beforeAll(async () => {
        const response = await axios.post(`${EVENT_URL}`, testEvent, {
            headers: { Authorization: `Bearer ${adminToken}` }
        });

        expect(response.status).toBe(200);
        expect(response.data).toHaveProperty("message", "Event created successfully");
        expect(response.data).toHaveProperty("event_id");

        eventID = response.data.event_id;
    });

    it("Get an event", async () => {
        const response = await axios.get(`${EVENT_URL}/`, {
            params: { eventID },
            headers: { Authorization: `Bearer ${userToken}` }
        });
        expect(response.status).toBe(200);
        expect(response.data).toHaveProperty("data");
    });

    it("Try to get an event with invalid eventID", async () => {
        try {
            await axios.get(`${EVENT_URL}`, {
                params: { eventID: "invalid-id" },
                headers: { Authorization: `Bearer ${userToken}` }
            });
        } catch (error: any) {
            expect(error.response?.status).toBe(405);
        }
    });

    it("Try to get an event without an ID", async () => {
        try {
            await axios.get(`${EVENT_URL}`, {
                headers: { Authorization: `Bearer ${userToken}` }
            });
        } catch (error: any) {
            expect(error.response?.status).toBe(405);
        }
    });

    it("Try to get an event without a token", async () => {
        try {
            await axios.get(`${EVENT_URL}`, {
                params: { eventID }
            });
        } catch (error: any) {
            expect(error.response?.status).toBe(405);
        }
    });
});


describe("Get all live events", () => {

    it("Get all live events", async () => {
        const response = await axios.get(`${EVENT_URL}/events/`, {
            params: { cursor: 0, limit: 10 },
            headers: { Authorization: `Bearer ${userToken}` }
        });
        expect(response.status).toBe(200);
        expect(response.data).toHaveProperty("data");
    });

    it("Try to get all live events without params, Get all live events with default", async () => {
        const response = await axios.get(`${EVENT_URL}/events/`, {
            params: { cursor: 0, limit: 10 },
            headers: { Authorization: `Bearer ${userToken}` }
        });
        expect(response.status).toBe(200);
        expect(response.data).toHaveProperty("data");
    });

    it("without a token: Unauthicated user try to get all live events", async () => {
        try {
            await axios.get(`${EVENT_URL}/events/`, {
                params: { cursor: 0, limit: 10 },
            });
        } catch (error: any) {
            expect(error.response?.status).toBe(403);
        }
    });
});