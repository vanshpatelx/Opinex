const queryList = {
    getUser: `SELECT id, email, password, type FROM "Auth" WHERE email = $1`
};

export default queryList;


// CREATE TABLE "Auth"(
//     "id" BIGINT NOT NULL,
//     "email" VARCHAR(255) NOT NULL,
//     "password" TEXT NOT NULL,
//     "type" VARCHAR(10) NOT NULL,
//     "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
//     "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
//     "inserted_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
// );
