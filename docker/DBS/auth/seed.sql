INSERT INTO "Auth" (id, email, password, type, created_at, updated_at, inserted_at) 
VALUES 
(1, 'test@gmail.com', '$2b$10$MoFnDOMosA2twi/9nObjbOgKvA.noI1cC0i9hIgfhz1yzKQ72vtmS', 
 'ADMIN', now(), now(), now()),
(2, 'testuser@gmail.com', '$2b$10$MoFnDOMosA2twi/9nObjbOgKvA.noI1cC0i9hIgfhz1yzKQ72vtmS', 
 'USER', now(), now(), now());
