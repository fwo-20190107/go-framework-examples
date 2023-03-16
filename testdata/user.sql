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
