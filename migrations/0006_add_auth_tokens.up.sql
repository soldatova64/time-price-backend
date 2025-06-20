CREATE TABLE auth_tokens
(
    "id"         SERIAL,
    "user_id"    INTEGER     NOT NULL,
    "token"      VARCHAR(32) NOT NULL UNIQUE,
    "created_at" TIMESTAMP   NOT NULL DEFAULT NOW(),
    "end_date"   TIMESTAMP   NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("user_id") REFERENCES "users"("id")
);

CREATE INDEX idx_auth_tokens_token ON auth_tokens(token);