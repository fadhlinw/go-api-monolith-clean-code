
-- +migrate Up
ALTER TABLE `users` ADD COLUMN `password` VARCHAR(255) NOT NULL AFTER `email`;

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `password`;
