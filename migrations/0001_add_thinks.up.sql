CREATE TABLE "thing"
(
    "id"         bigserial,
    "name"       varchar(128) NOT NULL,
    "pay_date"   date         NOT NULL DEFAULT CURRENT_DATE,
    "pay_price"  int          not null DEFAULT 0,
    "sale_date"  date                  DEFAULT null,
    "sale_price" bigint                DEFAULT null,
    "created_at" timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);