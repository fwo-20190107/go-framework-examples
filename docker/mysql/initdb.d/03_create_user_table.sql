USE `app`;

CREATE TABLE user (
    `user_id` int unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(16) NOT NULL,
    `authority` tinyint unsigned NOT NULL,
    PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='ユーザー情報';

INSERT INTO user (
    `user_id`,
    `name`,
    `authority`
) VALUES
(1, "admin", 99),
(2, "user", 20)
;
