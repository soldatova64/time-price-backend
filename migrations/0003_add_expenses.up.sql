CREATE TABLE "expense"
(
    "id"            BIGSERIAL    PRIMARY KEY,
    "thing_id"      BIGINT       NOT NULL,
    "sum"           BIGINT       NOT NULL,
    "description"   TEXT         NOT NULL,
    "created_at"    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "is_deleted"    BOOLEAN      NOT NULL DEFAULT FALSE,
    "deleted_at"    TIMESTAMP    WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY ("thing_id") REFERENCES "thing"("id")

);
CREATE INDEX idx_expense_deleted ON expense(is_deleted)