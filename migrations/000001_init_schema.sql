-- +migrate Up

CREATE TABLE `users` (`id` VARCHAR(128) PRIMARY KEY UNIQUE ,`first_name` VARCHAR(64), `last_name` VARCHAR(64), `nickname` VARCHAR(64), `password` VARCHAR(32), email VARCHAR(64), country VARCHAR(64), created_at DATETIME, updated_at DATETIME, UNIQUE(nickname), UNIQUE(email));

-- +migrate Down

DROP TABLE IF EXISTS employees;