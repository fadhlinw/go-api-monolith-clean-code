
-- +migrate Up
ALTER TABLE `users`
  ADD COLUMN `username` VARCHAR(255) NOT NULL AFTER `email`;

-- +migrate Down
ALTER TABLE `users`
  DROP COLUMN `username`;
