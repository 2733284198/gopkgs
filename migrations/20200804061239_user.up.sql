CREATE TABLE `users`(
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `email` VARCHAR(64) NOT NULL,
    `email_verified` TINYINT NOT NULL DEFAULT 0,
    `verification_token` VARCHAR(32) NULL,
    `hashed_password` VARCHAR(128) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_users_email` (`email`),
    UNIQUE KEY `idx_users_verification_token` (`verification_token`)
);
