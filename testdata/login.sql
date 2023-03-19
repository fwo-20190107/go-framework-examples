CREATE TABLE login (
    `login_id` TEXT NOT NULL PRIMARY KEY,
    `user_id` INTEGER NOT NULL,
    `last_signed_at` DATETIME,
    `password` TEXT NOT NULL
);

INSERT INTO login (
    `login_id`,
    `user_id`,
    `last_signed_at`,
    `password`
) VALUES
("admin", 1, null, "admin"),
("user", 2, null, "user")
;