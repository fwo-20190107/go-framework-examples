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