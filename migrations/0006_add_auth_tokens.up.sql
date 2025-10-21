CREATE TABLE auth_tokens
(
    "id"         BIGSERIAL,
    "user_id"    INTEGER     NOT NULL,
    "token"      VARCHAR(32) NOT NULL UNIQUE,
    "created_at" TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "end_date"   TIMESTAMP   NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("user_id") REFERENCES "users"("id")
);
