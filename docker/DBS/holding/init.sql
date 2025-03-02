CREATE TABLE "Holding" (
    "user_id" BIGINT NOT NULL PRIMARY KEY,
    "event_id" BIGINT NOT NULL,
    "price" INT NOT NULL,
    "quantity" INT NOT NULL,
    "lq" INT NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "inserted_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
);
