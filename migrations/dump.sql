INSERT INTO "users" ("id", "username", "email", "password")
VALUES
    (1, 'anna1', 'test@example.com', 'hashed_Qwerty123ffg'),
    (2,	'anna2', 'test@example.com',	'hashed_Qwerty123ffg'),
    (3,	'anna2', 'test@example.com',	'hashed_Qwerty123ffg')
;

INSERT INTO "thing" ("id", "name", "pay_date", "pay_price", "sale_date", "sale_price", "user_id")
VALUES
    (1, 'телефон', '2025-01-01', 50000, null, null, 1),
    (2, 'утюг', '2023-01-01', 5000, '2025-01-01', 3000, 1),
    (3, 'машина', '2020-01-01', 500000, null, null, 1)
;

INSERT INTO "expense" ("id", "thing_id", "sum", "description")
VALUES
    (1, 3, 5000, 'ремонт')
;

