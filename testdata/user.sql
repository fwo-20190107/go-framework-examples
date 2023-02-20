CREATE TABLE user (
    `user_id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    `name` TEXT NOT NULL,
    `authority` INTEGER NOT NULL
);

INSERT INTO user (
    `user_id`,
    `name`,
    `authority`
) VALUES
(1, "admin", 99),
(2, "user", 20)
;


CREATE TABLE login (
    `login_id` TEXT NOT NULL PRIMARY KEY,
    `user_id` INTEGER NOT NULL,
    `password` TEXT NOT NULL
);

INSERT INTO login (
    `login_id`,
    `user_id`,
    `password`
) VALUES
("admin", 1, "admin"),
("user", 2, "user")
;