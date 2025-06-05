CREATE TABLE "users"
(
    "id"         BIGSERIAL,
    "username"   VARCHAR(50)   NOT NULL,
    "email"      VARCHAR(100)  NOT NULL,
    "password"   VARCHAR(255)  NOT NULL,
    "created_at" TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "is_deleted" BOOLEAN       NOT NULL DEFAULT FALSE,
    "deleted_at" TIMESTAMP     WITH TIME ZONE DEFAULT NULL,
    PRIMARY KEY ("id")
);

CREATE INDEX idx_users_deleted ON users(is_deleted);