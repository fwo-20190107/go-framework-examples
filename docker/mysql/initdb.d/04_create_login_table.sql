USE `app`;

CREATE TABLE login (
    `login_id` varchar(16) NOT NULL,
    `user_id` int unsigned NOT NULL,
    `last_signed_at` datetime DEFAULT NULL,
    `password` varchar(32) NOT NULL,
    PRIMARY KEY (`login_id`),
    INDEX idx_login_01 (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='ログイン情報';

INSERT INTO login (
    `login_id`,
    `user_id`,
    `last_signed_at`,
    `password`
) VALUES
("admin", 1, null, "admin"),
("user", 2, null, "user")
;
