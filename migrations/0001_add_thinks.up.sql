CREATE TABLE "thing"
(
    "id"         BIGSERIAL,
    "name"       VARCHAR(128) NOT NULL,
    "pay_date"   DATE        NOT NULL DEFAULT CURRENT_DATE,
    "pay_price"  INT          NOT NULL DEFAULT 0,
    "sale_date"  DATE                  DEFAULT NULL,
    "sale_price" BIGINT                DEFAULT NULL,
    "created_at" TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);